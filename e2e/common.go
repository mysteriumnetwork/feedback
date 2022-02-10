package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/s3manager"
	"github.com/mysteriumnetwork/feedback/client"
	"github.com/stretchr/testify/assert"
)

// Error type
type Error struct {
	Message string  `json:"message"`
	Cause   *string `json:"Cause"`
}

// ErrorResponse type
type ErrorResponse struct {
	Errors []Error `json:"errors"`
	Code   int     `json:"code"`
}

// SendReportIntercomIssueRequest send an intercom issue request to the feedback API
func SendReportIntercomIssueRequest(t *testing.T, payload *client.CreateIntercomIssueRequest, filename string, fileContent []byte) *ErrorResponse {
	url := "http://localhost:8080/api/v1/intercom"

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	if len(fileContent) > 0 {
		part, err := writer.CreateFormFile("file", filename)
		assert.Nil(t, err)
		_, err = part.Write(fileContent)
		assert.Nil(t, err)
	}

	_ = writer.WriteField("userId", payload.UserId)
	_ = writer.WriteField("description", payload.Description)
	_ = writer.WriteField("email", payload.Email)
	_ = writer.WriteField("nodeIdentity", payload.NodeIdentity)
	_ = writer.WriteField("nodeCountry", payload.NodeCountry)
	_ = writer.WriteField("userType", payload.UserType)
	_ = writer.WriteField("ipType", payload.IpType)
	_ = writer.WriteField("ip", payload.Ip)
	_ = writer.Close()

	req, err := http.NewRequest("POST", url, body)
	assert.NoError(t, err)

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	var af struct{}
	return parseResp(t, resp, &af)
}

func parseResp(t *testing.T, resp *http.Response, obj interface{}) *ErrorResponse {
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var result ErrorResponse
		err := json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)
		result.Code = resp.StatusCode
		return &result
	}

	return nil
}

type s3Downloader struct {
	downloader *s3manager.Downloader
	client     *s3.Client
	bucket     string
}

func newS3Downloader(bucket string) (*s3Downloader, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, fmt.Errorf("could not load AWS configuration: %w", err)
	}
	cfg.EndpointResolver = aws.ResolveWithEndpointURL("http://localhost:9090")
	cfg.Region = endpoints.EuCentral1RegionID
	s3client := s3.New(cfg)
	s3client.ForcePathStyle = true
	downloader := &s3manager.Downloader{
		S3: s3client,
	}
	return &s3Downloader{
		downloader: downloader,
		client:     s3client,
		bucket:     bucket,
	}, nil
}

func (s3d *s3Downloader) getFileContent(t *testing.T, filename string) ([]byte, error) {
	paginator := s3.NewListObjectsV2Paginator(s3d.client.ListObjectsV2Request(&s3.ListObjectsV2Input{
		Bucket: &s3d.bucket,
	}))
	for paginator.Next(context.Background()) {
		page := paginator.CurrentPage()
		for _, obj := range page.Contents {
			fmt.Println(*obj.Key)
			if strings.Contains(*obj.Key, filename) {
				buf := aws.NewWriteAtBuffer([]byte{})
				_, err := s3d.downloader.Download(buf, &s3.GetObjectInput{
					Bucket: &s3d.bucket,
					Key:    obj.Key,
				})
				if err != nil {
					return nil, fmt.Errorf("download failed: %w", err)
				}
				return buf.Bytes(), nil
			}
		}
	}
	err := paginator.Err()
	if err != nil {
		return nil, fmt.Errorf("pagination error: %w", err)
	}
	return nil, fmt.Errorf("file not found")
}

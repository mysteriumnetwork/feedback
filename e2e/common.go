package e2e

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"testing"

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
func SendReportIntercomIssueRequest(t *testing.T, payload *client.CreateIntercomIssueRequest, fileContent []byte) *ErrorResponse {
	url := "http://localhost:8080/api/v1/intercom"

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	if len(fileContent) > 0 {
		part, err := writer.CreateFormFile("file", "example.txt")
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

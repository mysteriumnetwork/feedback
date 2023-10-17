package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/mysteriumnetwork/feedback/feedback"
	"github.com/mysteriumnetwork/feedback/infra/apierror"
)

const basePath = "/api/v1"

// FeedbackAPI client to submit user feedback
type FeedbackAPI struct {
	baseURL string
	http    *http.Client
}

// NewFeedbackAPI creates new FeedbackAPI
func NewFeedbackAPI(apiURL string) (*FeedbackAPI, error) {
	u, err := url.Parse(apiURL)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, basePath)

	return &FeedbackAPI{
		baseURL: u.String(),
		http: &http.Client{
			Timeout: 20 * time.Second,
		},
	}, nil
}

func (f *FeedbackAPI) apiURL(apiPath string) string {
	u, _ := url.Parse(f.baseURL) // Already validated in constructor
	u.Path = path.Join(u.Path, apiPath)
	return u.String()
}

// CreateGithubIssue creates a github issue
func (f *FeedbackAPI) CreateGithubIssue(request feedback.CreateGithubIssueRequest, logFilePath string) (*feedback.CreateGithubIssueResponse, *apierror.APIErrorResponse, error) {
	multipartReq := newMultipartRequest()
	err := multipartReq.addFileToMultipart("file", logFilePath)
	if err != nil {
		return nil, nil, fmt.Errorf("could not add file to multipart request: %w", err)
	}
	err = fillGithubFormFields(multipartReq, &request)
	if err != nil {
		return nil, nil, fmt.Errorf("could not fill form fields: %w", err)
	}

	body, contentType, err := multipartReq.finalize()
	if err != nil {
		return nil, nil, fmt.Errorf("could not finalize multipart request: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, f.apiURL("/github"), body)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create request: %w", err)
	}
	req.Header.Set("Content-Type", contentType)

	resp, err := f.http.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed request to feedback service: %w", err)
	}

	result := &feedback.CreateGithubIssueResponse{}
	apierror, err := parseResponse(resp, &result)
	if err != nil {
		return nil, nil, fmt.Errorf("could not parse response: %w", err)
	}

	return result, apierror, nil
}

// CreateIntercomIssue creates a intercom issue
func (f *FeedbackAPI) CreateIntercomIssue(request feedback.CreateIntercomIssueRequest, logFilePath string) (*feedback.CreateIntercomIssueResponse, *apierror.APIErrorResponse, error) {
	multipartReq := newMultipartRequest()
	err := multipartReq.addFileToMultipart("file", logFilePath)
	if err != nil {
		return nil, nil, fmt.Errorf("could not add file to multipart request: %w", err)
	}
	err = fillIntercomFormFields(multipartReq, &request)
	if err != nil {
		return nil, nil, fmt.Errorf("could not fill form fields: %w", err)
	}

	body, contentType, err := multipartReq.finalize()
	if err != nil {
		return nil, nil, fmt.Errorf("could not finalize multipart request: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, f.apiURL("/intercom"), body)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create request: %w", err)
	}
	req.Header.Set("Content-Type", contentType)

	resp, err := f.http.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed request to feedback service: %w", err)
	}

	result := &feedback.CreateIntercomIssueResponse{}
	apierror, err := parseResponse(resp, &result)
	if err != nil {
		return nil, nil, fmt.Errorf("could not parse response: %w", err)
	}

	return result, apierror, nil
}

// CreateBugReport creates a bug report
func (f *FeedbackAPI) CreateBugReport(request feedback.CreateBugReportRequest, logFilePath string) (*feedback.CreateBugReportResponse, *apierror.APIErrorResponse, error) {
	multipartReq := newMultipartRequest()
	err := multipartReq.addFileToMultipart("file", logFilePath)
	if err != nil {
		return nil, nil, fmt.Errorf("could not add file to multipart request: %w", err)
	}
	err = fillBugReportFormFields(multipartReq, &request)
	if err != nil {
		return nil, nil, fmt.Errorf("could not fill form fields: %w", err)
	}

	body, contentType, err := multipartReq.finalize()
	if err != nil {
		return nil, nil, fmt.Errorf("could not finalize multipart request: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, f.apiURL("/bug-report"), body)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create request: %w", err)
	}
	req.Header.Set("Content-Type", contentType)

	resp, err := f.http.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed request to feedback service: %w", err)
	}

	result := &feedback.CreateBugReportResponse{}
	apierror, err := parseResponse(resp, &result)
	if err != nil {
		return nil, nil, fmt.Errorf("could not parse response: %w", err)
	}

	return result, apierror, nil
}

func fillBugReportFormFields(multipartRequest *multipartRequest, request *feedback.CreateBugReportRequest) error {
	err := multipartRequest.addFieldToMultipart("nodeIdentity", request.NodeIdentity)
	if err != nil {
		return fmt.Errorf("could not add nodeIdentity to multipart request: %w", err)
	}
	err = multipartRequest.addFieldToMultipart("email", request.Email)
	if err != nil {
		return fmt.Errorf("could not add email to multipart request: %w", err)
	}
	err = multipartRequest.addFieldToMultipart("description", request.Description)
	if err != nil {
		return fmt.Errorf("could not add description to multipart request: %w", err)
	}
	return nil
}

func fillIntercomFormFields(multipartRequest *multipartRequest, request *feedback.CreateIntercomIssueRequest) error {
	err := multipartRequest.addFieldToMultipart("userId", request.UserId)
	if err != nil {
		return fmt.Errorf("could not add userId to multipart request: %w", err)
	}
	err = multipartRequest.addFieldToMultipart("nodeIdentity", request.NodeIdentity)
	if err != nil {
		return fmt.Errorf("could not add nodeIdentity to multipart request: %w", err)
	}
	err = multipartRequest.addFieldToMultipart("userType", request.UserType)
	if err != nil {
		return fmt.Errorf("could not add userType to multipart request: %w", err)
	}
	err = multipartRequest.addFieldToMultipart("nodeCountry", request.NodeCountry)
	if err != nil {
		return fmt.Errorf("could not add nodeCountry to multipart request: %w", err)
	}
	err = multipartRequest.addFieldToMultipart("ipType", request.IpType)
	if err != nil {
		return fmt.Errorf("could not add ipType to multipart request: %w", err)
	}
	err = multipartRequest.addFieldToMultipart("ip", request.Ip)
	if err != nil {
		return fmt.Errorf("could not add ip to multipart request: %w", err)
	}
	err = multipartRequest.addFieldToMultipart("description", request.Description)
	if err != nil {
		return fmt.Errorf("could not add description to multipart request: %w", err)
	}
	err = multipartRequest.addFieldToMultipart("email", request.Email)
	if err != nil {
		return fmt.Errorf("could not add email to multipart request: %w", err)
	}
	return nil
}

func fillGithubFormFields(multipartRequest *multipartRequest, request *feedback.CreateGithubIssueRequest) error {
	err := multipartRequest.addFieldToMultipart("userId", request.UserId)
	if err != nil {
		return fmt.Errorf("could not add userId to multipart request: %w", err)
	}
	err = multipartRequest.addFieldToMultipart("description", request.Description)
	if err != nil {
		return fmt.Errorf("could not add description to multipart request: %w", err)
	}
	err = multipartRequest.addFieldToMultipart("email", request.Email)
	if err != nil {
		return fmt.Errorf("could not add email to multipart request: %w", err)
	}
	return nil
}

type multipartRequest struct {
	buffer *bytes.Buffer
	writer *multipart.Writer
}

func newMultipartRequest() *multipartRequest {
	buffer := new(bytes.Buffer)
	return &multipartRequest{
		buffer: buffer,
		writer: multipart.NewWriter(buffer),
	}
}

func (mr *multipartRequest) addFileToMultipart(fieldname, filename string) error {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("could not read input file: %w", err)
	}
	part, err := mr.writer.CreateFormFile("file", filename)
	if err != nil {
		return fmt.Errorf("could not add file to request: %w", err)
	}
	_, err = part.Write(fileContent)
	if err != nil {
		return fmt.Errorf("could not write file part: %w", err)
	}
	return nil
}

func (mr *multipartRequest) addFieldToMultipart(fieldname, value string) error {
	return mr.writer.WriteField(fieldname, value)
}

func (mr *multipartRequest) finalize() (*bytes.Buffer, string, error) {
	err := mr.writer.Close()
	if err != nil {
		return nil, "", fmt.Errorf("could not close multipart writer: %w", err)
	}
	return mr.buffer, mr.writer.FormDataContentType(), nil
}

func parseResponse(resp *http.Response, v any) (*apierror.APIErrorResponse, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could read response: %w", err)
	}
	if resp.StatusCode >= 400 {
		apierror := apierror.APIErrorResponse{}
		err = json.Unmarshal(body, &apierror)
		if err != nil {
			return nil, fmt.Errorf("could not parse error response: %w", err)
		}
		return &apierror, nil
	}
	err = json.Unmarshal(body, v)
	if err != nil {
		return nil, fmt.Errorf("could not parse response: %w", err)
	}
	return nil, nil
}

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
func (f *FeedbackAPI) CreateGithubIssue(request CreateGithubIssueRequest) (response *CreateGithubIssueResult, err error) {
	multipartReq, err := newCreateGithubIssueRequest(f.apiURL("/github"), request)
	if err != nil {
		return nil, fmt.Errorf("could not create multipart request: %w", err)
	}

	resp, err := f.http.Do(multipartReq)
	if err != nil {
		return nil, fmt.Errorf("failed request to feedback service: %w", err)
	}
	return parseCreateGithubIssueResult(resp)
}

// CreateIntercomIssue creates a intercom issue
func (f *FeedbackAPI) CreateIntercomIssue(request CreateIntercomIssueRequest) (response *CreateIntercomIssueResult, err error) {
	multipartReq, err := newCreateIntercomIssueRequest(f.apiURL("/intercom"), request)
	if err != nil {
		return nil, fmt.Errorf("could not create multipart request: %w", err)
	}

	resp, err := f.http.Do(multipartReq)
	if err != nil {
		return nil, fmt.Errorf("failed request to feedback service: %w", err)
	}
	return parseCreateIntercomIssueResult(resp)
}

func newCreateGithubIssueRequest(uri string, req CreateGithubIssueRequest) (multipartReq *http.Request, err error) {
	fileContent, err := os.ReadFile(req.Filepath)
	if err != nil {
		return nil, fmt.Errorf("could not read input file: %w", err)
	}
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", req.Filepath)
	if err != nil {
		return nil, fmt.Errorf("could not add file to request: %w", err)
	}
	_, err = part.Write(fileContent)
	if err != nil {
		return nil, fmt.Errorf("could not write file part: %w", err)
	}

	_ = writer.WriteField("userId", req.UserId)
	_ = writer.WriteField("description", req.Description)
	_ = writer.WriteField("email", req.Email)
	_ = writer.Close()

	request, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())
	return request, err
}

func newCreateIntercomIssueRequest(uri string, req CreateIntercomIssueRequest) (multipartReq *http.Request, err error) {
	fileContent, err := os.ReadFile(req.Filepath)
	if err != nil {
		return nil, fmt.Errorf("could not read input file: %w", err)
	}
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", req.Filepath)
	if err != nil {
		return nil, fmt.Errorf("could not add file to request: %w", err)
	}
	_, err = part.Write(fileContent)
	if err != nil {
		return nil, fmt.Errorf("could not write file part: %w", err)
	}

	_ = writer.WriteField("userId", req.UserId)
	_ = writer.WriteField("description", req.Description)
	_ = writer.WriteField("email", req.Email)
	_ = writer.WriteField("nodeIdentity", req.NodeIdentity)
	_ = writer.WriteField("nodeCountry", req.NodeCountry)
	_ = writer.WriteField("userType", req.UserType)
	_ = writer.WriteField("ipType", req.IpType)
	_ = writer.WriteField("ip", req.Ip)
	_ = writer.Close()

	request, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())
	return request, err
}

// CreateGithubIssueRequest create github issue request
type CreateGithubIssueRequest struct {
	UserId      string `json:"userId"`
	Description string `json:"description"`
	Email       string `json:"email"`
	Filepath    string `json:"file"`
}

// CreateIntercomIssueRequest create intercom issue request
type CreateIntercomIssueRequest struct {
	UserId       string `json:"userId"`
	Description  string `json:"description"`
	Email        string `json:"email"`
	Filepath     string `json:"file"`
	NodeIdentity string `json:"nodeIdentity"`
	UserType     string `json:"userType"`
	NodeCountry  string `json:"nodeCountry"`
	IpType       string `json:"ipType"`
	Ip           string `json:"ip"`
}

// CreateGithubIssueResult represents create github issue response (either successful or error)
type CreateGithubIssueResult struct {
	Response     *GithubIssueCreated
	Errors       *ErrorResponse
	HTTPResponse *http.Response
	Success      bool
}

// CreateIntercomIssueResult represents create intercom issue response (either successful or error)
type CreateIntercomIssueResult struct {
	Errors       *ErrorResponse
	HTTPResponse *http.Response
	Success      bool
}

// GithubIssueCreated represents successful response (issue created)
type GithubIssueCreated struct {
	IssueId string `json:"issueId"`
}

func parseCreateGithubIssueResult(httpRes *http.Response) (*CreateGithubIssueResult, error) {
	res := &CreateGithubIssueResult{HTTPResponse: httpRes}
	resJSON, err := io.ReadAll(httpRes.Body)
	if err != nil {
		res.Success = false
		res.Errors = singleErrorResponse("could not parse feedback response: " + err.Error())
		return res, err
	}
	if res.HTTPResponse.StatusCode >= 400 {
		res.Success = false
		res.Errors = &ErrorResponse{[]Error{}}
		err = json.Unmarshal(resJSON, res.Errors)
		if err != nil {
			res.Errors = singleErrorResponse("could not parse feedback response: " + err.Error())
			return res, err
		}
	} else {
		res.Response = &GithubIssueCreated{}
		err := json.Unmarshal(resJSON, res.Response)
		if err != nil {
			res.Success = false
			res.Errors = singleErrorResponse("could not parse feedback response: " + err.Error())
			return res, err
		}
		res.Success = true
	}
	return res, nil
}

func parseCreateIntercomIssueResult(httpRes *http.Response) (*CreateIntercomIssueResult, error) {
	res := &CreateIntercomIssueResult{HTTPResponse: httpRes}
	resJSON, err := io.ReadAll(httpRes.Body)
	if err != nil {
		res.Success = false
		res.Errors = singleErrorResponse("could not parse feedback response: " + err.Error())
		return res, err
	}
	if res.HTTPResponse.StatusCode >= 400 {
		res.Success = false
		res.Errors = &ErrorResponse{[]Error{}}
		err = json.Unmarshal(resJSON, res.Errors)
		if err != nil {
			res.Errors = singleErrorResponse("could not parse feedback response: " + err.Error() + ". Response was: " + string(resJSON))
			return res, fmt.Errorf("parsing \"%s\" failed: %w", string(resJSON), err)
		}
	} else {
		res.Success = true
	}
	return res, nil
}

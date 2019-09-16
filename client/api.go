/*
 * Copyright (C) 2019 The "MysteriumNetwork/feedback" Authors.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/pkg/errors"
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
		// For now not using fmt.Errorf with %w for compatibility with go <1.13
		return nil, errors.Wrap(err, "could not create multipart request")
	}

	resp, err := f.http.Do(multipartReq)
	if err != nil {
		return nil, errors.Wrap(err, "failed request to feedback service")
	}
	return parseCreateGithubIssueResult(resp)
}

func newCreateGithubIssueRequest(uri string, req CreateGithubIssueRequest) (multipartReq *http.Request, err error) {
	fileContent, err := ioutil.ReadFile(req.Filepath)
	if err != nil {
		return nil, errors.Wrap(err, "could not read input file")
	}
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", req.Filepath)
	if err != nil {
		return nil, errors.Wrap(err, "could not add file to request")
	}
	_, err = part.Write(fileContent)
	if err != nil {
		return nil, errors.Wrap(err, "could not write file part")
	}

	_ = writer.WriteField("userId", req.UserId)
	_ = writer.WriteField("description", req.Description)
	_ = writer.WriteField("email", req.Email)
	_ = writer.Close()

	request, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, errors.Wrap(err, "could not create request")
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

// CreateGithubIssueResult represents create github issue response (either successful or error)
type CreateGithubIssueResult struct {
	Response     *GithubIssueCreated
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
	resJSON, err := ioutil.ReadAll(httpRes.Body)
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

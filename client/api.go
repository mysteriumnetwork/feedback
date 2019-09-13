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

	"github.com/mysteriumnetwork/feedback/infra"
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
func (f *FeedbackAPI) CreateGithubIssue(request CreateGithubIssueRequest) (response *CreateGithubIssueResponse, err error) {
	multipartReq, err := newCreateGithubIssueRequest(f.apiURL("/github"), request)
	if err != nil {
		// For now not using fmt.Errorf with %w for compatibility with go <1.13
		return nil, errors.Wrap(err, "could not create multipart request")
	}

	resp, err := f.http.Do(multipartReq)
	if err != nil {
		return nil, errors.Wrap(err, "failed request to feedback service")
	}
	return parseCreateGithubIssueResponse(resp)
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
	return http.NewRequest("POST", uri, body)
}

// CreateGithubIssueRequest create github issue request
type CreateGithubIssueRequest struct {
	UserId      string `json:"userId"`
	Description string `json:"description"`
	Email       string `json:"email"`
	Filepath    string `json:"file"`
}

// CreateGithubIssueResponse represents create github issue response (either successful or error)
type CreateGithubIssueResponse struct {
	Result       *CreateGithubIssueResponse
	Errors       *infra.ErrorResponse
	HTTPResponse *http.Response
	Success      bool
}

func parseCreateGithubIssueResponse(httpRes *http.Response) (*CreateGithubIssueResponse, error) {
	res := &CreateGithubIssueResponse{
		HTTPResponse: httpRes,
		Success:      true,
	}
	resJSON, err := ioutil.ReadAll(httpRes.Body)
	if res.HTTPResponse.StatusCode < 400 {
		err := json.Unmarshal(resJSON, res.Result)
		if err != nil {
			return res, err
		}
	} else {
		res.Success = false
		if err != nil {
			return res, err
		}
		err = json.Unmarshal(resJSON, res.Errors)
		if err != nil {
			return res, err
		}
	}
	return res, nil
}

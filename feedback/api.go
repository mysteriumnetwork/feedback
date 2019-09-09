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

package feedback

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"

	log "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"github.com/mysteriumnetwork/feedback/apierr"
	"github.com/mysteriumnetwork/feedback/storage"
	"github.com/mysteriumnetwork/feedback/target/github"
)

// Endpoint user feedback endpoint
type Endpoint struct {
	reporter *github.Reporter
	storage  *storage.Storage
}

// NewEndpoint creates new Endpoint
func NewEndpoint(reporter *github.Reporter, storage *storage.Storage) *Endpoint {
	return &Endpoint{reporter: reporter, storage: storage}
}

// CreateGithubIssueRequest create github issue request
type CreateGithubIssueRequest struct {
	UserId      string
	Description string
	Email       string
	File        *multipart.FileHeader
}

// CreateGithubIssueResponse represents a successful github issue creation
type CreateGithubIssueResponse struct {
	IssueId string `json:"issueId"`
}

// ParseGithubIssueRequest parses CreateGithubIssueRequest from HTTP request
func ParseGithubIssueRequest(c *gin.Context) (form CreateGithubIssueRequest, errors []error) {
	var ok bool
	form.UserId, ok = c.GetPostForm("userId")
	if !ok {
		errors = append(errors, apierr.Required("userId"))
	}

	form.Description, ok = c.GetPostForm("description")
	if !ok {
		errors = append(errors, apierr.Required("description"))
	}

	form.Email, ok = c.GetPostForm("email")
	if !ok {
		errors = append(errors, apierr.Required("email"))
	}

	var err error
	form.File, err = c.FormFile("file")
	if err != nil {
		errors = append(errors, apierr.Required("file"))
	}

	return form, errors
}

// CreateGithubIssue creates a new Github issue with user report
func (e *Endpoint) CreateGithubIssue(c *gin.Context) {
	form, requestErrs := ParseGithubIssueRequest(c)
	if len(requestErrs) > 0 {
		c.JSON(http.StatusBadRequest, apierr.Multiple(requestErrs))
		return
	}

	tempFile, err := ioutil.TempFile("", form.File.Filename)
	if err != nil {
		err := fmt.Errorf("could not allocate a temporary file: %w", err)
		c.JSON(http.StatusInternalServerError, apierr.Single(err))
		return
	}
	defer func() {
		err := os.Remove(tempFile.Name())
		if err != nil {
			_ = log.Warn("Failed to remove temp file", err)
		}
	}()

	err = c.SaveUploadedFile(form.File, tempFile.Name())
	if err != nil {
		err := fmt.Errorf("could not save the uploaded file: %w", err)
		c.JSON(http.StatusInternalServerError, apierr.Single(err))
		return
	}

	logURL, err := e.storage.Upload(tempFile.Name())
	if err != nil {
		err := fmt.Errorf("could not upload file to the storage: %w", err)
		c.JSON(http.StatusServiceUnavailable, apierr.Single(err))
		return
	}

	issueId, err := e.reporter.ReportIssue(&github.Report{
		UserId:      form.UserId,
		Description: form.Description,
		Email:       form.Email,
		LogURL:      *logURL,
	})
	if err != nil {
		err := fmt.Errorf("could not report issue to github: %w", err)
		c.JSON(http.StatusServiceUnavailable, apierr.Single(err))
		return
	}

	log.Infof("Created github issue %q from request %+v", issueId, form)
	c.JSON(http.StatusOK, &CreateGithubIssueResponse{
		IssueId: issueId,
	})
}

// RegisterRoutes registers feedback API routes
func (e *Endpoint) RegisterRoutes(r gin.IRoutes) {
	r.POST("/github", e.CreateGithubIssue)
}

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

	"github.com/gin-gonic/gin"
	"github.com/mysteriumnetwork/feedback/apierr"
	"github.com/pkg/errors"
)

// Endpoint user feedback endpoint
type Endpoint struct {
}

// NewEndpoint creates new Endpoint
func NewEndpoint() *Endpoint {
	return &Endpoint{}
}

// GithubIssueForm create github issue request
type GithubIssueForm struct {
	UserId      string
	Description string
	Email       string
	File        *multipart.FileHeader
}

// ParseGithubIssueForm parses GithubIssueForm from HTTP request
func ParseGithubIssueForm(c *gin.Context) (form GithubIssueForm, errors []error) {
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
	form, formErrors := ParseGithubIssueForm(c)
	if len(formErrors) > 0 {
		c.JSON(http.StatusBadRequest, apierr.Multiple(formErrors))
		return
	}

	tempFile, err := ioutil.TempFile("", form.File.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apierr.Single(errors.Wrap(err, "could not allocate a temporary file")))
		return
	}

	err = c.SaveUploadedFile(form.File, tempFile.Name())
	if err != nil {
		c.JSON(http.StatusInternalServerError, apierr.Single(errors.Wrap(err, "could not save the uploaded file")))
	}

	email := c.PostForm("email")
	fmt.Println("Email: " + email)
	description := c.PostForm("description")
	fmt.Println("Description: " + description)
	fmt.Println("Uploaded: " + tempFile.Name())

	c.JSON(http.StatusOK, "file uploaded")
}

// RegisterRoutes registers feedback API routes
func (e *Endpoint) RegisterRoutes(r gin.IRoutes) {
	r.POST("/github", e.CreateGithubIssue)
}

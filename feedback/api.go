package feedback

import (
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"

	log "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"github.com/mysteriumnetwork/feedback/infra"
	"github.com/mysteriumnetwork/feedback/infra/apierror"
)

// Endpoint user feedback endpoint
type Endpoint struct {
	githubReporter   Reporter
	intercomReporter Reporter
	storage          Uploader
	rateLimiter      *infra.RateLimiter
	skipFileUpload   *bool
}

// Uploader uploads a file and returns the url
type Uploader interface {
	Upload(filepath string) (url *url.URL, err error)
}

// Report represents an issue report
type Report struct {
	UserId       string
	NodeIdentity string
	UserType     string
	NodeCountry  string
	IpType       string
	Ip           string
	Description  string
	Email        string
	LogURL       url.URL
}

// Reporter reports issues and returns their id
type Reporter interface {
	ReportIssue(report *Report) (issueId string, err error)
}

// NewEndpoint creates new Endpoint
func NewEndpoint(githubReporter Reporter, intercomReporter Reporter, storage Uploader, rateLimiter *infra.RateLimiter, skipFileUpload *bool) *Endpoint {
	return &Endpoint{githubReporter: githubReporter, storage: storage, rateLimiter: rateLimiter, intercomReporter: intercomReporter, skipFileUpload: skipFileUpload}
}

// CreateGithubIssueRequest create github issue request
// swagger:parameters createGithubIssue
type CreateGithubIssueRequest struct {
	// in: formData
	// required: true
	UserId string `json:"userId"`
	// in: formData
	// required: true
	Description string `json:"description"`
	// in: formData
	Email string `json:"email"`
	// in: formData
	// required: true
	// swagger:file
	File *multipart.FileHeader `json:"file"`
}

// CreateGithubIssueResponse represents a successful github issue creation
// swagger:model
type CreateGithubIssueResponse struct {
	IssueId string `json:"issueId"`
}

// ParseGithubIssueRequest parses CreateGithubIssueRequest from HTTP request
func ParseGithubIssueRequest(c *gin.Context) (form CreateGithubIssueRequest, errors []error) {
	_, err := c.MultipartForm()
	if err != nil {
		errors = append(errors, apierror.APIError{Message: "could not parse form: " + err.Error()})
		return form, errors
	}

	var ok bool
	form.UserId, ok = c.GetPostForm("userId")
	if !ok {
		errors = append(errors, apierror.Required("userId"))
	}

	form.Description, ok = c.GetPostForm("description")
	if !ok {
		errors = append(errors, apierror.Required("description"))
	}

	form.Email = c.PostForm("email")

	form.File, err = c.FormFile("file")
	if err != nil {
		errors = append(errors, apierror.Required("file"))
	}

	return form, errors
}

// CreateGithubIssue creates a new Github issue with user report
//
// swagger:operation POST /github createGithubIssue
// ---
// summary: Creates a new Github issue with user report
// description: 1 request per minute is allowed
//
// produces:
// - application/json
// consumes:
// - multipart/form-data
// responses:
//   '200':
//     description: Issue created in Github
//     schema:
//       "$ref": "#/definitions/CreateGithubIssueResponse"
//   '400':
//     description: Bad request
//     schema:
//       "$ref": "#/definitions/APIErrorResponse"
//   '429':
//     description: Too many requests
//   '500':
//     description: Internal server error
//     schema:
//       "$ref": "#/definitions/APIErrorResponse"
//
func (e *Endpoint) CreateGithubIssue(c *gin.Context) {
	form, requestErrs := ParseGithubIssueRequest(c)
	if len(requestErrs) > 0 {
		c.JSON(http.StatusBadRequest, apierror.Multiple(requestErrs))
		return
	}

	tempFile, err := ioutil.TempFile("", path.Base(form.File.Filename))
	if err != nil {
		apiError := apierror.New("Could not allocate a temporary file", err)
		_ = log.Error(apiError.Wrapped())
		c.JSON(http.StatusInternalServerError, apiError.ToResponse())
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
		apiError := apierror.New("Could not save the uploaded file", err)
		_ = log.Error(apiError.Wrapped())
		c.JSON(http.StatusInternalServerError, apiError.ToResponse())
		return
	}

	var logURL *url.URL
	if !*e.skipFileUpload {
		logURL, err = e.storage.Upload(tempFile.Name())
		if err != nil {
			apiError := apierror.New("could not upload file to the storage", err)
			_ = log.Error(apiError.Wrapped())
			c.JSON(http.StatusServiceUnavailable, apiError.ToResponse())
			return
		}
	}

	report := &Report{
		UserId:      form.UserId,
		Description: form.Description,
		Email:       form.Email,
	}

	if logURL != nil {
		report.LogURL = *logURL
	}

	issueId, err := e.githubReporter.ReportIssue(report)
	if err != nil {
		apiError := apierror.New("could not report issue to github", err)
		_ = log.Error(apiError.Wrapped())
		c.JSON(http.StatusServiceUnavailable, apiError.ToResponse())
		return
	}

	log.Infof("Created github issue %q from request %+v", issueId, form)
	c.JSON(http.StatusOK, &CreateGithubIssueResponse{
		IssueId: issueId,
	})
}

// CreateIntercomIssueRequest create intercom issue request
// swagger:parameters createIntercomIssue
type CreateIntercomIssueRequest struct {
	// in: formData
	// required: false
	UserId string `json:"userId"`
	// in: formData
	// required: true
	NodeIdentity string `json:"nodeIdentity"`
	// in: formData
	// required: false
	UserType string `json:"userType"`
	// in: formData
	// required: false
	NodeCountry string `json:"nodeCountry"`
	// in: formData
	// required: false
	IpType string `json:"ipType"`
	// in: formData
	// required: false
	Ip string `json:"ip"`
	// in: formData
	// required: true
	Description string `json:"description"`
	// in: formData
	// required: false
	Email string `json:"email"`
	// in: formData
	// required: true
	// swagger:file
	File *multipart.FileHeader `json:"file"`
}

// CreateIntercomIssueResponse represents a successful intercom issue creation
// swagger:model
type CreateIntercomIssueResponse struct {
	ConversationId string `json:"conversationId"`
}

// ParseIntercomIssueRequest parses CreateIntercomIssueRequest from HTTP request
func ParseIntercomIssueRequest(c *gin.Context) (form CreateIntercomIssueRequest, errors []error) {
	_, err := c.MultipartForm()
	if err != nil {
		errors = append(errors, apierror.APIError{Message: "could not parse form: " + err.Error()})
		return form, errors
	}

	form.UserId = c.PostForm("userId")

	var ok bool
	form.NodeIdentity, ok = c.GetPostForm("nodeIdentity")
	if !ok || form.NodeIdentity == "" {
		errors = append(errors, apierror.Required("nodeIdentity"))
	}

	form.UserType = c.PostForm("userType")
	form.NodeCountry = c.PostForm("nodeCountry")
	form.IpType = c.PostForm("ipType")
	form.Ip = c.PostForm("ip")

	form.Description, ok = c.GetPostForm("description")
	if !ok {
		errors = append(errors, apierror.Required("description"))
	}
	if len(form.Description) < 30 {
		errors = append(errors, apierror.Custom("description", "too short"))
	}

	form.Email = c.PostForm("email")
	if form.Email == "" && form.UserId == "" {
		errors = append(errors, apierror.Required("email or userId"))
	}

	form.File, err = c.FormFile("file")
	if err != nil {
		errors = append(errors, apierror.Required("file"))
	}

	return form, errors
}

// CreateIntercomIssue creates a new Intercom chat with the issue with user report
//
// swagger:operation POST /intercom createIntercomIssue
// ---
// summary: Creates a new Intercom chat with the issue with user report
// description: 1 request per minute is allowed
//
// produces:
// - application/json
// consumes:
// - multipart/form-data
// responses:
//   '200':
//     description: Chat created in Intercom
//     schema:
//       "$ref": "#/definitions/CreateIntercomIssueResponse"
//   '400':
//     description: Bad request
//     schema:
//       "$ref": "#/definitions/APIErrorResponse"
//   '429':
//     description: Too many requests
//   '500':
//     description: Internal server error
//     schema:
//       "$ref": "#/definitions/APIErrorResponse"
//
func (e *Endpoint) CreateIntercomIssue(c *gin.Context) {
	form, requestErrs := ParseIntercomIssueRequest(c)
	if len(requestErrs) > 0 {
		c.JSON(http.StatusBadRequest, apierror.Multiple(requestErrs))
		return
	}

	tempFile, err := ioutil.TempFile("", path.Base(form.File.Filename))
	if err != nil {
		apiError := apierror.New("Could not allocate a temporary file", err)
		_ = log.Error(apiError.Wrapped())
		c.JSON(http.StatusInternalServerError, apiError.ToResponse())
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
		apiError := apierror.New("Could not save the uploaded file", err)
		_ = log.Error(apiError.Wrapped())
		c.JSON(http.StatusInternalServerError, apiError.ToResponse())
		return
	}

	var logURL *url.URL
	if !*e.skipFileUpload {
		logURL, err = e.storage.Upload(tempFile.Name())
		if err != nil {
			apiError := apierror.New("could not upload file to the storage", err)
			_ = log.Error(apiError.Wrapped())
			c.JSON(http.StatusServiceUnavailable, apiError.ToResponse())
			return
		}
	}

	report := &Report{
		UserId:       form.UserId,
		NodeIdentity: form.NodeIdentity,
		UserType:     form.UserType,
		NodeCountry:  form.NodeCountry,
		IpType:       form.IpType,
		Ip:           form.Ip,
		Description:  form.Description,
		Email:        form.Email,
	}

	if logURL != nil {
		report.LogURL = *logURL
	}

	conversationId, err := e.intercomReporter.ReportIssue(report)
	if err != nil {
		apiError := apierror.New("could not create intercom conversation", err)
		_ = log.Error(apiError.Wrapped())
		c.JSON(http.StatusServiceUnavailable, apiError.ToResponse())
		return
	}

	log.Infof("Created intercom conversation from request %+v", form)
	c.JSON(http.StatusOK, &CreateIntercomIssueResponse{
		ConversationId: conversationId,
	})
}

// RegisterRoutes registers feedback API routes
func (e *Endpoint) RegisterRoutes(r gin.IRoutes) {
	r.POST("/github", e.rateLimiter.Handler(), e.CreateGithubIssue)
	r.POST("/intercom", e.rateLimiter.Handler(), e.CreateIntercomIssue)
}

package feedback

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/mysteriumnetwork/feedback/infra"
	"github.com/mysteriumnetwork/feedback/infra/apierror"
	"github.com/rs/zerolog/log"
)

// Endpoint user feedback endpoint
type Endpoint struct {
	githubReporter   Reporter
	intercomReporter ReporterWithMessage
	storage          Uploader
	rateLimiter      *infra.RateLimiter
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

// ReporterWithMessage reports issues and generates bug report messages
type ReporterWithMessage interface {
	Reporter
	GetBugReportMessage(report *Report) (message string, err error)
}

// NewEndpoint creates new Endpoint
func NewEndpoint(githubReporter Reporter, intercomReporter ReporterWithMessage, storage Uploader, rateLimiter *infra.RateLimiter) *Endpoint {
	return &Endpoint{githubReporter: githubReporter, storage: storage, rateLimiter: rateLimiter, intercomReporter: intercomReporter}
}

// CreateGithubIssueRequest create github issue request
type CreateGithubIssueRequest struct {
	UserId      string                `form:"userId"`
	Description string                `form:"description"`
	Email       string                `form:"email"`
	File        *multipart.FileHeader `form:"file"`
}

// CreateGithubIssueResponse represents a successful github issue creation
type CreateGithubIssueResponse struct {
	IssueId string `json:"issueId"`
}

// ParseGithubIssueRequest parses CreateGithubIssueRequest from HTTP request
func ParseGithubIssueRequest(c *gin.Context) (CreateGithubIssueRequest, []error) {
	errors := make([]error, 0)
	form := CreateGithubIssueRequest{}

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

	var err error
	form.File, err = c.FormFile("file")
	if err != nil {
		errors = append(errors, apierror.Required("file"))
	}

	return form, errors
}

// CreateGithubIssue creates a new Github issue with user report
//
// Deprecated: use CreateBugReport instead
func (e *Endpoint) CreateGithubIssue(c *gin.Context) {
	form, requestErrs := ParseGithubIssueRequest(c)
	if len(requestErrs) > 0 {
		c.JSON(http.StatusBadRequest, apierror.Multiple(requestErrs))
		return
	}

	logURL, err := e.uploadFile(c, form.File)
	if err != nil {
		log.Error().Err(err).Msg("could not upload file")
		c.JSON(http.StatusServiceUnavailable, apierror.New("could not upload file", err).ToResponse())
		return
	}

	report := &Report{
		NodeIdentity: form.UserId,
		Description:  form.Description,
		Email:        form.Email,
		LogURL:       *logURL,
	}

	issueId, err := e.githubReporter.ReportIssue(report)
	if err != nil {
		log.Error().Err(err).Msg("could not report issue to github")
		c.JSON(http.StatusServiceUnavailable, apierror.New("could not report issue to github", err).ToResponse())
		return
	}

	log.Info().Str("issue_id", issueId).Interface("form", form).Msg("Created github issue from request")
	c.JSON(http.StatusOK, &CreateGithubIssueResponse{
		IssueId: issueId,
	})
}

// CreateIntercomIssueRequest create intercom issue request
type CreateIntercomIssueRequest struct {
	UserId       string `form:"userId" validate:"optional"`
	NodeIdentity string `form:"nodeIdentity" validate:"required"`
	UserType     string `form:"userType" validate:"optional"`
	NodeCountry  string `form:"nodeCountry" validate:"optional"`
	IpType       string `form:"ipType" validate:"optional"`
	Ip           string `form:"ip" validate:"optional"`
	Description  string `form:"description" validate:"required" minLength:"30"`
	Email        string `form:"email" validate:"required"`
	File         *multipart.FileHeader
}

// CreateIntercomIssueResponse represents a successful intercom issue creation
type CreateIntercomIssueResponse struct {
	ConversationId string `json:"conversationId"`
}

// ParseIntercomIssueRequest parses CreateIntercomIssueRequest from HTTP request
func ParseIntercomIssueRequest(c *gin.Context) (CreateIntercomIssueRequest, []error) {
	errors := make([]error, 0)
	form := CreateIntercomIssueRequest{}

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
	if !ok || form.Description == "" {
		errors = append(errors, apierror.Required("description"))
	} else if len(form.Description) < 30 {
		errors = append(errors, apierror.Custom("description", "too short"))
	}

	form.Email = c.PostForm("email")
	if form.Email == "" && form.UserId == "" {
		errors = append(errors, apierror.Required("email or userId"))
	}

	var err error
	form.File, err = c.FormFile("file")
	if err != nil {
		errors = append(errors, apierror.Required("file"))
	}

	return form, errors
}

// CreateIntercomIssue creates a new Intercom chat with the issue with user report
// Deprecated: use CreateBugReport instead
//
// @Tags bug-report
// @Summary create a new bug in intercom
// @Param request formData feedback.CreateIntercomIssueRequest true "bug report request"
// @Param file formData file true "log file"
// @Accepts multipart/form-data
// @Produces json
// @Success 200 {object} feedback.CreateIntercomIssueResponse
// @Failure 400 {object} apierror.APIErrorResponse
// @Failure 429
// @Failure 500 {object} apierror.APIErrorResponse
// @Router /v1/intercom [post]
// @Deprecated
func (e *Endpoint) CreateIntercomIssue(c *gin.Context) {
	form, requestErrs := ParseIntercomIssueRequest(c)
	if len(requestErrs) > 0 {
		c.JSON(http.StatusBadRequest, apierror.Multiple(requestErrs))
		return
	}

	logURL, err := e.uploadFile(c, form.File)
	if err != nil {
		log.Error().Err(err).Msg("could not upload file")
		c.JSON(http.StatusServiceUnavailable, apierror.New("could not upload file", err).ToResponse())
		return
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
		LogURL:       *logURL,
	}

	conversationId, err := e.intercomReporter.ReportIssue(report)
	if err != nil {
		log.Error().Err(err).Msg("could not create intercom conversation")
		c.JSON(http.StatusServiceUnavailable, apierror.New("could not create intercom conversation", err).ToResponse())
		return
	}
	log.Info().Interface("form", form).Msg("Created intercom conversation from request")

	_, err = e.githubReporter.ReportIssue(report)
	if err != nil {
		log.Error().Err(err).Msg("could not report issue to github")
		//this is added for conveniance for the support team, we will not alert users
		//that it failed as intercom report was successful
	}

	c.JSON(http.StatusOK, &CreateIntercomIssueResponse{
		ConversationId: conversationId,
	})
}

// CreateBugReportRequest create bug report request
type CreateBugReportRequest struct {
	NodeIdentity string `form:"nodeIdentity" example:"0xF0345F6251Bef9447A08766b9DA2B07b28aD80B0" validate:"required"`
	Description  string `form:"description" minLength:"30" validate:"required"`
	Email        string `form:"email" validate:"required"`
	File         *multipart.FileHeader
}

// ParseBugReportRequest parses CreateBugReportRequest from HTTP request
func ParseBugReportRequest(c *gin.Context) (CreateBugReportRequest, []error) {
	errors := make([]error, 0)
	form := CreateBugReportRequest{}

	var ok bool
	form.NodeIdentity, ok = c.GetPostForm("nodeIdentity")
	if !ok || form.NodeIdentity == "" {
		errors = append(errors, apierror.Required("nodeIdentity"))
	} else if !common.IsHexAddress(form.NodeIdentity) {
		errors = append(errors, apierror.Custom("nodeIdentity", "invalid address"))
	}

	form.Description, ok = c.GetPostForm("description")
	if !ok || form.Description == "" {
		errors = append(errors, apierror.Required("description"))
	} else if len(form.Description) < 30 {
		errors = append(errors, apierror.Custom("description", "too short"))
	}

	form.Email, ok = c.GetPostForm("email")
	if !ok || form.Email == "" {
		errors = append(errors, apierror.Required("email"))
	}

	var err error
	form.File, err = c.FormFile("file")
	if err != nil {
		errors = append(errors, apierror.Required("file"))
	}

	return form, errors
}

// CreateBugReportResponse represents a successful bug report creation
type CreateBugReportResponse struct {
	IssueId      string `json:"issueId"`
	Message      string `json:"message"`
	Email        string `json:"email"`
	NodeIdentity string `json:"node_identity"`
}

// CreateBugReport creates a new bug report.
//
// @Tags bug-report
// @Summary create a new bug report
// @Param request formData feedback.CreateBugReportRequest true "bug report request"
// @Param file formData file true "log file"
// @Accepts multipart/form-data
// @Produces json
// @Success 200 {object} feedback.CreateBugReportResponse
// @Failure 400 {object} apierror.APIErrorResponse
// @Failure 429
// @Failure 500 {object} apierror.APIErrorResponse
// @Router /v1/bug-report [post]
func (e *Endpoint) CreateBugReport(c *gin.Context) {
	form, requestErrs := ParseBugReportRequest(c)
	if len(requestErrs) > 0 {
		c.JSON(http.StatusBadRequest, apierror.Multiple(requestErrs))
		return
	}

	logURL, err := e.uploadFile(c, form.File)
	if err != nil {
		log.Error().Err(err).Msg("could not upload file")
		c.JSON(http.StatusServiceUnavailable, apierror.New("could not upload file", err).ToResponse())
		return
	}

	report := &Report{
		NodeIdentity: form.NodeIdentity,
		Description:  form.Description,
		Email:        form.Email,
		LogURL:       *logURL,
	}

	message, err := e.intercomReporter.GetBugReportMessage(report)
	if err != nil {
		log.Error().Err(err).Msg("could not get bug report message")
		c.JSON(http.StatusServiceUnavailable, apierror.New("could not get bug report message", err).ToResponse())
		return
	}

	issueId, err := e.githubReporter.ReportIssue(report)
	if err != nil {
		log.Error().Err(err).Msg("could not report issue to github")
		//this is added for conveniance for the support team, we will not alert users that it failed
	}

	c.JSON(http.StatusOK, &CreateBugReportResponse{
		IssueId:      issueId,
		Message:      message,
		Email:        form.Email,
		NodeIdentity: form.NodeIdentity,
	})
}

// RegisterRoutes registers feedback API routes
func (e *Endpoint) RegisterRoutes(r gin.IRoutes) {
	r.POST("/github", e.rateLimiter.Handler(), e.CreateGithubIssue)
	r.POST("/intercom", e.rateLimiter.Handler(), e.CreateIntercomIssue)
	r.POST("/bug-report", e.rateLimiter.Handler(), e.CreateBugReport)
}

func (e *Endpoint) uploadFile(c *gin.Context, file *multipart.FileHeader) (*url.URL, error) {
	tempFile, err := os.CreateTemp("", path.Base(file.Filename))
	if err != nil {
		return nil, fmt.Errorf("could not allocate a temporary file: %w", err)
	}
	defer func() {
		err := os.Remove(tempFile.Name())
		if err != nil {
			log.Warn().Err(err).Msg("failed to remove temp file")
		}
	}()

	err = c.SaveUploadedFile(file, tempFile.Name())
	if err != nil {
		return nil, fmt.Errorf("could not save the uploaded file: %w", err)
	}

	logURL, err := e.storage.Upload(tempFile.Name())
	if err != nil {
		return nil, fmt.Errorf("could not upload file to the storage: %w", err)
	}

	return logURL, nil
}

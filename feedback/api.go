package feedback

import (
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"

	log "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"github.com/mysteriumnetwork/feedback/infra"
	"github.com/mysteriumnetwork/feedback/infra/apierror"
)

// Endpoint user feedback endpoint
type Endpoint struct {
	reporter    *Reporter
	storage     *Storage
	rateLimiter *infra.RateLimiter
}

// NewEndpoint creates new Endpoint
func NewEndpoint(reporter *Reporter, storage *Storage, rateLimiter *infra.RateLimiter) *Endpoint {
	return &Endpoint{reporter: reporter, storage: storage, rateLimiter: rateLimiter}
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

	logURL, err := e.storage.Upload(tempFile.Name())
	if err != nil {
		apiError := apierror.New("could not upload file to the storage", err)
		_ = log.Error(apiError.Wrapped())
		c.JSON(http.StatusServiceUnavailable, apiError.ToResponse())
		return
	}

	issueId, err := e.reporter.ReportIssue(&Report{
		UserId:      form.UserId,
		Description: form.Description,
		Email:       form.Email,
		LogURL:      *logURL,
	})
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

// RegisterRoutes registers feedback API routes
func (e *Endpoint) RegisterRoutes(r gin.IRoutes) {
	r.POST("/github", e.rateLimiter.Handler(), e.CreateGithubIssue)
}

package apierror

import "fmt"

// APIErrorResponse represent HTTP error payload
type APIErrorResponse struct {
	Errors []APIError `json:"errors"`
}

// APIError represents a single error in an APIErrorResponse
type APIError struct {
	Message string `json:"message"`
	Cause   error  `json:"-"`
}

// New creates a new APIError
func New(msg string, cause error) APIError {
	return APIError{
		Message: msg,
		Cause:   cause,
	}
}

// NewMsg creates a new APIError without a cause
func NewMsg(msg string) APIError {
	return APIError{
		Message: msg,
	}
}

// Error represents APIError as string
func (e APIError) Error() string {
	return e.Message
}

// Wrapped returns APIError with a cause wrapped
func (e APIError) Wrapped() error {
	return fmt.Errorf("%s: %w", e.Message, e.Cause)
}

// ToResponse returns APIError to APIErrorResponse containing a single APIError
func (e APIError) ToResponse() *APIErrorResponse {
	return Multiple([]error{e})
}

// Single creates an error response containing a single error
func Single(apiError APIError) *APIErrorResponse {
	return Multiple([]error{apiError})
}

// Multiple creates an error response containing multiple errors
func Multiple(errors []error) *APIErrorResponse {
	var apiErrors []APIError
	for _, err := range errors {
		apiErrors = append(apiErrors, APIError{Message: err.Error()})
	}
	return &APIErrorResponse{Errors: apiErrors}
}

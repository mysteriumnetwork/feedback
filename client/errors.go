package client

// ErrorResponse represent HTTP error payload
// swagger:model
type ErrorResponse struct {
	Errors []Error `json:"errors"`
}

// Error represents a single error in an ErrorResponse
type Error struct {
	Message string `json:"message"`
}

func singleErrorResponse(msg string) *ErrorResponse {
	return &ErrorResponse{[]Error{
		{Message: msg},
	}}
}

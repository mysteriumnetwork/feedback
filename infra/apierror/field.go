package apierror

// Required creates a new RequiredFieldError
func Required(field string) RequiredFieldError {
	return RequiredFieldError{field: field}
}

// RequiredFieldError represents a missing required field in HTTP request
type RequiredFieldError struct {
	field string
}

// Error returns user friendly message for RequiredFieldError
func (f RequiredFieldError) Error() string {
	return "field is required: " + f.field
}

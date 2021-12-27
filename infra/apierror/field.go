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

// Custom creates a new CustomFieldError
func Custom(field string, customReason string) CustomFieldError {
	return CustomFieldError{field: field, customReason: customReason}
}

// CustomFieldError represents a custom error for a field in HTTP request
type CustomFieldError struct {
	field        string
	customReason string
}

// Error returns user friendly message for CustomFieldError
func (f CustomFieldError) Error() string {
	return "field is " + f.customReason + ": " + f.field
}

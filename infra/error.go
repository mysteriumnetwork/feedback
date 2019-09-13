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

package infra

// ErrorResponse represent HTTP error payload
// swagger:model
type ErrorResponse struct {
	Errors []Error `json:"errors"`
}

// Error represents a single error in an ErrorResponse
type Error struct {
	Message string `json:"message"`
}

// Single creates an error response containing a single error
func Single(err error) *ErrorResponse {
	return Multiple([]error{err})
}

// Multiple creates an error response containing multiple errors
func Multiple(errors []error) *ErrorResponse {
	var apierrs []Error
	for _, err := range errors {
		apierrs = append(apierrs, Error{Message: err.Error()})
	}
	return &ErrorResponse{Errors: apierrs}
}

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

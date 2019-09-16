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

package apierror

import "fmt"

// APIErrorResponse represent HTTP error payload
// swagger:model
type APIErrorResponse struct {
	Errors []APIError `json:"errors"`
}

// APIError represents a single error in an APIErrorResponse
type APIError struct {
	Message string `json:"message"`
	Cause   error
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

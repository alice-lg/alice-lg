package http

// Improve error handling
// Create api.ErrorResponses based on errors returned from server.
// Strip out potentially sensitive information, eg. connection errors
// to internal IP addresses.

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/alice-lg/alice-lg/pkg/api"
)

// ErrResourceNotFoundError is a 404 error
type ErrResourceNotFoundError struct{}

// Error implements the error interface and returns
// the error message
func (err *ErrResourceNotFoundError) Error() string {
	return "resource not found"
}

// Variables
var (
	ErrSourceNotFound = &ErrResourceNotFoundError{}
)

// Error tags
const (
	TagGenericError      = "GENERIC_ERROR"
	TagConnectionRefused = "CONNECTION_REFUSED"
	TagConnectionTimeout = "CONNECTION_TIMEOUT"
	TagResourceNotFound  = "NOT_FOUND"
	TagValidationError   = "VALIDATION_ERROR"
)

// Error codes
const (
	CodeGeneric           = 42
	CodeConnectionRefused = 100
	CodeConnectionTimeout = 101
	CodeValidationError   = 400
	CodeResourceNotFound  = 404
)

// Error status codes
const (
	StatusError            = http.StatusInternalServerError
	StatusResourceNotFound = http.StatusNotFound
	StatusValidationError  = http.StatusBadRequest
)

// Handle an error and create a error API response
func apiErrorResponse(
	routeserverID string,
	err error,
) (api.ErrorResponse, int) {
	code := CodeGeneric
	message := err.Error()
	tag := TagGenericError
	status := StatusError

	switch e := err.(type) {
	case *ErrResourceNotFoundError:
		tag = TagResourceNotFound
		code = CodeResourceNotFound
		status = StatusResourceNotFound
	case *url.Error:
		if strings.Contains(message, "connection refused") {
			tag = TagConnectionRefused
			code = CodeConnectionRefused
			message = "Connection refused while dialing the API"
		} else if e.Timeout() {
			tag = TagConnectionTimeout
			code = CodeConnectionTimeout
			message = "Connection timed out when connecting to the backend API"
		}
	}

	switch err {
	case ErrQueryTooShort:
		tag = TagValidationError
		code = CodeValidationError
		status = StatusValidationError
		message = "the query is too short"
	case ErrQueryIncomplete:
		tag = TagValidationError
		code = CodeValidationError
		status = StatusValidationError
		message = "the query is incomplete"
	}

	return api.ErrorResponse{
		Code:          code,
		Tag:           tag,
		Message:       message,
		RouteserverID: routeserverID,
	}, status
}

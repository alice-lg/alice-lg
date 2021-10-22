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

// ResourceNotFoundError is a 404 error
type ResourceNotFoundError struct{}

// Error implements the error interface and returns
// the error message
func (err *ResourceNotFoundError) Error() string {
	return "resource not found"
}

// Variables
var (
	SOURCE_NOT_FOUND_ERROR = &ResourceNotFoundError{}
)

// Error Constants
const (
	GENERIC_ERROR_TAG      = "GENERIC_ERROR"
	CONNECTION_REFUSED_TAG = "CONNECTION_REFUSED"
	CONNECTION_TIMEOUT_TAG = "CONNECTION_TIMEOUT"
	RESOURCE_NOT_FOUND_TAG = "NOT_FOUND"
)

const (
	GENERIC_ERROR_CODE      = 42
	CONNECTION_REFUSED_CODE = 100
	CONNECTION_TIMEOUT_CODE = 101
	RESOURCE_NOT_FOUND_CODE = 404
)

const (
	ERROR_STATUS              = http.StatusInternalServerError
	RESOURCE_NOT_FOUND_STATUS = http.StatusNotFound
)

// Handle an error and create a error API response
func apiErrorResponse(routeserverId string, err error) (api.ErrorResponse, int) {
	code := GENERIC_ERROR_CODE
	message := err.Error()
	tag := GENERIC_ERROR_TAG
	status := ERROR_STATUS

	switch e := err.(type) {
	case *ResourceNotFoundError:
		tag = RESOURCE_NOT_FOUND_TAG
		code = RESOURCE_NOT_FOUND_CODE
		status = RESOURCE_NOT_FOUND_STATUS
	case *url.Error:
		if strings.Contains(message, "connection refused") {
			tag = CONNECTION_REFUSED_TAG
			code = CONNECTION_REFUSED_CODE
			message = "Connection refused while dialing the API"
		} else if e.Timeout() {
			tag = CONNECTION_TIMEOUT_TAG
			code = CONNECTION_TIMEOUT_CODE
			message = "Connection timed out when connecting to the backend API"
		}
	}

	return api.ErrorResponse{
		Code:          code,
		Tag:           tag,
		Message:       message,
		RouteserverId: routeserverId,
	}, status
}

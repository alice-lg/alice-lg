package main

// Improve error handling
// Create api.ErrorResponses based on errors returned from server.
// Strip out potentially sensitive information, eg. connection errors
// to internal IP addresses.

import (
	"net/url"
	"strings"

	"github.com/alice-lg/alice-lg/backend/api"
)

const (
	GENERIC_ERROR_TAG      = "GENERIC_ERROR"
	CONNECTION_REFUSED_TAG = "CONNECTION_REFUSED"
	CONNECTION_TIMEOUT_TAG = "CONNECTION_TIMEOUT"
)

const (
	GENERIC_ERROR_CODE      = 42
	CONNECTION_REFUSED_CODE = 100
	CONNECTION_TIMEOUT_CODE = 101
)

func apiErrorResponse(routeserverId string, err error) api.ErrorResponse {
	code := GENERIC_ERROR_CODE
	message := err.Error()
	tag := GENERIC_ERROR_TAG

	switch e := err.(type) {
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
	}
}

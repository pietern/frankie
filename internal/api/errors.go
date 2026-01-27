package api

import "errors"

var (
	// ErrAuthRequired is returned when authentication is required but not provided
	ErrAuthRequired = errors.New("authentication required")

	// ErrInvalidCredentials is returned when login credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrForbidden is returned when access is forbidden
	ErrForbidden = errors.New("forbidden")

	// ErrNotFound is returned when a resource is not found
	ErrNotFound = errors.New("not found")

	// ErrBadRequest is returned when the request is malformed
	ErrBadRequest = errors.New("bad request")

	// ErrServerError is returned when the server returns an error
	ErrServerError = errors.New("server error")

	// ErrNetwork is returned when there's a network error
	ErrNetwork = errors.New("network error")

	// ErrSmartTradingNotEnabled is returned when smart trading is not enabled
	ErrSmartTradingNotEnabled = errors.New("smart trading is not enabled for this user")

	// ErrSmartChargingNotEnabled is returned when smart charging is not enabled
	ErrSmartChargingNotEnabled = errors.New("smart charging is not enabled for this user")

	// ErrNotSupportedInCountry is returned when the request is not supported in the user's country
	ErrNotSupportedInCountry = errors.New("request not supported in this country")
)

// APIError represents an error from the Frank Energie API
type APIError struct {
	Message string
	Path    []string
	Code    string
}

func (e *APIError) Error() string {
	return e.Message
}

// GraphQLError represents a GraphQL error response
type GraphQLError struct {
	Message    string                 `json:"message"`
	Path       []interface{}          `json:"path,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

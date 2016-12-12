package elserror

// APIError defines an error that should be thrown by the REST API
type APIError struct {
	Code       uint16
	Message    string
	StatusCode int
}

var (
	// GeneralError applies to an unrecoverable internal error
	GeneralError = &APIError{Code: 1, Message: "General Error", StatusCode: 500}
)

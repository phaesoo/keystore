package resp

// Code type
type Code = string

// General code
const (
	CodeUnknownError    Code = "UNKNOWN_ERROR"
	CodeInvalidRequest  Code = "INVALID_REQUEST"
	CodeInternalError   Code = "INTERNAL_ERROR"
	CodeValidationError Code = "VALIDATION_ERROR"
	CodeUnknownResource Code = "UNKNOWN_RESOURCE"
	CodeAuthRequired    Code = "AUTH_REQUIRED"
	CodeUnknownEndpoint Code = "UNKNOWN_ENDPOINT"
)

// Business error code
const (
	CodeVerificationFailed Code = "VERIFICATION_FAILED"
	CodeInvalidTokenState  Code = "INVALID_TOKEN_STATE"
)

// Status type
type Status = string

// Status
const (
	StatusSuccess Status = "success"
	StatusFail    Status = "fail"
	StatusError   Status = "error"
)

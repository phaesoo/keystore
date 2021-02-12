package resp

import (
	"net/http"
)

// Success responds with status 200 and the provided data and message
// status code >= 200 < 400
func Success(w http.ResponseWriter, statusCode int, data interface{}, message string) {
	WriteJSON(w, NewSuccess(data, message), statusCode)
}

// Fail responds with the given status code and the provided code and message.
// Code should be machine readable. Data can optionally be provided
// status code >= 400 < 500
func Fail(w http.ResponseWriter, status int, data interface{}, message string, code string) {
	WriteJSON(w, NewFail(data, code, message), status)
}

// Error responds with 501 and the provided code and error message
// status code >= 500
func Error(w http.ResponseWriter, err error, code string) {
	// Default error code
	if code == "" {
		code = CodeUnknownError
	}
	WriteJSON(w, NewError(code, err.Error()), http.StatusInternalServerError)
}

// OK is wrapper func for Success with status code 200
func OK(w http.ResponseWriter, data interface{}, message string) {
	Success(w, http.StatusOK, data, message)
}

// InvalidRequest is used to notify that the request is invalid and cannot be handled
func InvalidRequest(w http.ResponseWriter) {
	Fail(w, http.StatusUnprocessableEntity, nil, "Invalid request body.", CodeInvalidRequest)
}

// UnprocessableError is used to return error with 422 HTTP status code.
func UnprocessableError(w http.ResponseWriter, message, code string) {
	if message == "" {
		message = "UnprocessableError"
	}
	Fail(w, http.StatusUnprocessableEntity, nil, message, code)
}

// ValidationError is used to notify that the request's parameters/body/form-data contained invalid data
func ValidationError(w http.ResponseWriter, code string) {
	if code == "" {
		code = CodeValidationError
	}
	Fail(w, http.StatusBadRequest, nil, "Request did not pass validation checks.", code)
}

// UnknownResource is used when the URL is valid, but the requested resource could not be found
func UnknownResource(w http.ResponseWriter) {
	Fail(w, http.StatusNotFound, nil, "The requested resource could not be found.", CodeUnknownResource)
}

// Unauthorized is used when the user does not have a valid token
func Unauthorized(w http.ResponseWriter) {
	Fail(w, http.StatusUnauthorized, nil, "Request is not authorized.", CodeAuthRequired)
}

// BadRequest
func BadRequest(w http.ResponseWriter, code string) {
	Fail(w, http.StatusBadRequest, nil, "Bad request.", code)
}

// Package resp provides convenience methods for handling JSON responses
package resp

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// Resp is the standard response format
type Resp struct {
	Status  Status      `json:"status"`         // Response status class
	Code    Code        `json:"code,omitempty"` // A machine-readable error code in the case of failure/error responses.
	Data    interface{} `json:"data"`           // A data payload provided for success/failure responses.
	Message string      `json:"message"`        // A user-friendly message detailing the result of the operation.
}

// NewSuccess creates a successful JSON response with the given data and optional message.
func NewSuccess(data interface{}, message string) Resp {
	if data == nil {
		data = struct{}{}
	}
	return Resp{
		Status:  StatusSuccess,
		Data:    data,
		Message: message,
	}
}

// NewFail creates a failure JSON response with the given data, code, and optional message.
func NewFail(data interface{}, code string, message string) Resp {
	if data == nil {
		data = struct{}{}
	}
	return Resp{
		Status:  StatusFail,
		Data:    data,
		Code:    code,
		Message: message,
	}
}

// NewError creates a new error response with the given code and optional message.
func NewError(code string, message string) Resp {
	return Resp{
		Status:  StatusError,
		Code:    code,
		Message: message,
	}
}

// WriteJSON writes out the given data to the `http.ResponseWriter`, along with required headers.
func WriteJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	setContentJSON(w)
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
}

// Get extracts a `resp.Resp` shaped response payload from a `http.Response`. The data will
// be available in the `dataPtr`, and the other JSON fields are returned.
func Get(body io.Reader, dataPtr interface{}) (Resp, error) {
	response := Resp{
		Data: dataPtr,
	}
	if err := json.NewDecoder(body).Decode(&response); err != nil {
		return Resp{}, errors.Wrap(err, "malformed JSON payload")
	}
	return response, nil
}

func setContentJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

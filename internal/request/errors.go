package request

import "errors"

// request errors
var (
	ErrInvalidContentType = errors.New("Invalid content type")
	ErrExpectedEOF        = errors.New("Unexpected trailing data exists")
)

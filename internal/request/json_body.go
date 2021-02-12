package request

import (
	"encoding/json"
	"net/http"

	"io"
)

const contentType = "application/json"

// BindJSON binds request body into interface v
func BindJSON(r *http.Request, v interface{}) error {
	if err := validateContentType(r, contentType); err != nil {
		return err
	}
	if err := decodeBody(r, v); err != nil {
		return err
	}
	return nil
}

func decodeBody(r *http.Request, v interface{}) error {
	body := r.Body
	defer body.Close()

	// decode body content
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(v); err != nil {
		return err
	}

	// ensure there is not trailing data in the request body
	if err := ensureEOF(decoder.Buffered()); err != nil {
		return err
	}
	if err := ensureEOF(body); err != nil {
		return err
	}
	return nil
}

func ensureEOF(r io.Reader) error {
	switch n, err := r.Read([]byte{0}); err {
	case io.EOF:
		if n == 0 {
			return nil
		}
		fallthrough
	case nil:
		return ErrExpectedEOF
	default:
		return nil
	}
} 

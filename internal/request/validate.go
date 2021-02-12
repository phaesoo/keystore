package request

import (
	"mime"
	"net/http"
)

func validateContentType(r *http.Request, contentType string) error {
	cType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil && err != mime.ErrInvalidMediaParameter {
		return err
	}
	if cType != contentType {
		return ErrInvalidContentType
	}
	return nil
}

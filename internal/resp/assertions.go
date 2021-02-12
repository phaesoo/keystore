package resp

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// AssertFail confirms that the response is a failure
func AssertFail(t *testing.T, res *httptest.ResponseRecorder, code, message string) {
	t.Helper()
	Assert(t, res, 0, StatusFail, message, code)
}

// AssertOK confirms that the response was successful
func AssertOK(t *testing.T, res *httptest.ResponseRecorder, message string) {
	t.Helper()
	Assert(t, res, http.StatusOK, StatusSuccess, message, "")
}

// AssertError confirms that the response was an error
func AssertError(t *testing.T, res *httptest.ResponseRecorder, code, message string) {
	t.Helper()
	Assert(t, res, http.StatusInternalServerError, StatusError, message, code)
}

// Assert makes an assertion about a Response, with optional status, message, and code checks.
func Assert(t *testing.T, resRec *httptest.ResponseRecorder, statusCode int, status, message, code string) {
	t.Helper()

	bodyBytes, _ := ioutil.ReadAll(resRec.Body)
	resRec.Body = bytes.NewBuffer(bodyBytes)

	var res Resp
	err := json.Unmarshal(bodyBytes, &res)

	asserts := require.New(t)

	asserts.Contains(resRec.Header().Get("content-type"), "application/json")
	asserts.NoError(err, string(bodyBytes), "error unmarshaling json")

	asserts.Equal(status, res.Status, "unexpected status for body: %s", string(bodyBytes))

	if statusCode != 0 {
		asserts.Equal(statusCode, resRec.Code, "unexpected code for body: %s", string(bodyBytes))
	}

	if message != "" {
		asserts.Equal(message, res.Message, "unexpected message")
	}

	if code != "" {
		asserts.Equal(code, res.Code, "unexpected code")
	}
}

// ExtractData extracts the `data` payload inside of the JSON response
func ExtractData(t *testing.T, res *httptest.ResponseRecorder, data interface{}) {
	t.Helper()
	_, err := Get(res.Body, data)
	require.NoError(t, err, "error while unmarshaling json")
}

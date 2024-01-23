package middleware

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/devdammit/shekel/pkg/log"
)

func copyBody(r *http.Request, logger log.Logger) []byte {
	buffer := new(bytes.Buffer)

	// Copy body
	_, _ = buffer.ReadFrom(r.Body)
	err := r.Body.Close()
	if err != nil {
		logger.With(log.Err(err)).Warn("error closing request body")
	}

	// Init new body from buffer
	r.Body = ioutil.NopCloser(buffer)

	return buffer.Bytes()
}

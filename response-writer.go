// Package imageopti Bypass default response writer to override headers and body
package imageopti

import (
	"bytes"
	"fmt"
	"net/http"
)

type responseWriter struct {
	buffer       bytes.Buffer
	bypassHeader bool
	wroteHeader  bool // Control when to write header

	http.ResponseWriter
}

func (r *responseWriter) WriteHeader(statusCode int) {
	if !r.bypassHeader {
		r.ResponseWriter.WriteHeader(statusCode)
	}
}

func (r *responseWriter) Write(p []byte) (int, error) {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}

	i, err := r.buffer.Write(p)
	if err != nil {
		return i, fmt.Errorf("unable to write response body: %w", err)
	}

	return i, nil
}

func (r *responseWriter) Flush() {
	if flusher, ok := r.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

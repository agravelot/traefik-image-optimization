package image_optimizer

import (
	"bytes"
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
	return r.buffer.Write(p)
}

func (r *responseWriter) Flush() {
	if flusher, ok := r.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

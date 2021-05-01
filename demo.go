// Package plugindemo a demo plugin.
package plugindemo

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/agravelot/plugindemo/optimizer"
)

// Config the plugin configuration.
type Config struct {
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// Demo a Demo plugin.
type Demo struct {
	next http.Handler
	name string
}

// New created a new Demo plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {

	log.Println("Loading image optimization plugin...")
	// TODO Init and config

	return &Demo{
		next: next,
		name: name,
	}, nil
}

func (a *Demo) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// Ignore non image requests

	if IsImageRequest(req.Header.Get("accept")) {
		// Return cahed result here.
	}

	wrappedWriter := &responseWriter{
		ResponseWriter: rw,
	}

	a.next.ServeHTTP(wrappedWriter, req)

	bodyBytes := wrappedWriter.buffer.Bytes()

	if !IsImageResponse(rw.Header().Get("content-type")) {
		return
	}

	// Delegates the Content-Length Header creation to the final body write.
	rw.Header().Del("Content-Length")

	optimizer, err := optimizer.New("imaginary")
	if err != nil {
		panic(err)
	}

	optimized, err := optimizer.Optimize(bodyBytes, "", "", 75)
	if err != nil {
		panic(err)
	}

	if _, err := rw.Write(optimized); err != nil {
		log.Printf("unable to write rewrited body: %v", err)
	}

}

// IsImageResponse Determine with Content-Type header if the response is an image.
func IsImageResponse(contentType string) bool {
	return strings.HasPrefix(contentType, "image/")
}

// IsImageRequest Determine with Accept header if the request ask for an image.
func IsImageRequest(acceptHeader string) bool {
	if len(acceptHeader) == 0 {
		return false
	}

	accepts := strings.Split(acceptHeader, ",")

	for _, value := range accepts {
		println(value)
		// If start with "text/html", return false
		if value == "text/html" {
			return false
		}

		// Avoid unseless loop
		if strings.HasPrefix(value, "image/") {
			return true
		}
	}

	return true
}

package image_optimizer

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/agravelot/image_optimizer/config"
	"github.com/agravelot/image_optimizer/processor"
)

// Config the plugin configuration.
type Config struct {
	config.Config
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// Demo a Demo plugin.
type Demo struct {
	config Config
	next   http.Handler
	name   string
}

// New created a new Demo plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {

	log.Println("Loading image optimization plugin...")
	// TODO Init and config
	fmt.Printf("config : %+v\n", config)

	if config.Processor == "" {
		return nil, fmt.Errorf("processor must be defined")
	}

	return &Demo{
		config: *config,
		next:   next,
		name:   name,
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

	processor, err := processor.New(a.config.Config)
	if err != nil {
		panic(err)
	}

	optimized, err := processor.Optimize(bodyBytes, "", "", 75)
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

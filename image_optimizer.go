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

const (
	contentLength = "Content-Length"
	contentType   = "Content-Type"
)

func (a *Demo) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// Ignore non image requests

	// TODO Check if cacheable
	// Return cached result here.

	wrappedWriter := &responseWriter{
		ResponseWriter: rw,
	}

	a.next.ServeHTTP(wrappedWriter, req)

	bodyBytes := wrappedWriter.buffer.Bytes()

	if !IsImageResponse(rw.Header().Get(contentType)) {
		return
	}

	// Delegates the Content-Length Header creation to the final body write.
	rw.Header().Del(contentLength)

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

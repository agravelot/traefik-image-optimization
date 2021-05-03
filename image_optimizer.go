package image_optimizer

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/agravelot/image_optimizer/cache"
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

// ImageOptimizer middleware plugin base.
type ImageOptimizer struct {
	config Config
	next   http.Handler
	name   string
}

// New created a new ImageOptimizer plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {

	log.Println("Loading image optimization plugin...")
	// TODO Init and config
	fmt.Printf("config : %+v\n", config)

	if config.Processor == "" {
		return nil, fmt.Errorf("processor must be defined")
	}

	return &ImageOptimizer{
		config: *config,
		next:   next,
		name:   name,
	}, nil
}

const (
	contentLength = "Content-Length"
	contentType   = "Content-Type"
)

func (a *ImageOptimizer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// TODO Check if cacheable

	// Return cached result here.
	c, err := cache.New(a.config.Config)
	if err != nil {
		panic(err)
	}

	key, err := cache.Tokenize(req)
	if err != nil {
		panic(err)
	}

	if v, err := c.Get(key); err == nil {
		_, err = rw.Write(v)
		if err != nil {
			panic(err)
		}
		return
	}

	wrappedWriter := &responseWriter{
		ResponseWriter: rw,
	}

	a.next.ServeHTTP(wrappedWriter, req)

	bodyBytes := wrappedWriter.buffer.Bytes()

	if !isImageResponse(rw) {
		_, err = rw.Write(bodyBytes)
		if err != nil {
			panic(err)
		}
		return
	}

	// Delegates the Content-Length Header creation to the final body write.
	rw.Header().Del(contentLength)

	p, err := processor.New(a.config.Config)
	if err != nil {
		panic(err)
	}

	optimized, err := p.Optimize(bodyBytes, "", "", 75)
	if err != nil {
		panic(err)
	}

	_, err = rw.Write(optimized)
	if err != nil {
		log.Printf("unable to write rewrited body: %v", err)
		panic(err)
	}

	err = c.Set(key, optimized, 100*time.Second)
	if err != nil {
		panic(err)
	}
}

// isImageResponse Determine with Content-Type header if the response is an image.
func isImageResponse(rw http.ResponseWriter) bool {
	return strings.HasPrefix(rw.Header().Get(contentType), "image/")
}

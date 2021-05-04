package image_optimizer

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
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

// ImageOptimizer middleware plugin struct.
type ImageOptimizer struct {
	next http.Handler
	name string
	p    processor.Processor
	c    cache.Cache
}

// New created a new ImageOptimizer plugin.
func New(ctx context.Context, next http.Handler, conf *Config, name string) (http.Handler, error) {

	log.Println("Loading image optimization plugin...")
	fmt.Printf("config : %+v\n", conf)

	if conf.Processor == "" {
		return nil, fmt.Errorf("processor must be defined")
	}

	c, err := cache.New(conf.Config)
	if err != nil {
		panic(err)
	}

	p, err := processor.New(conf.Config)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	return &ImageOptimizer{
		p:    p,
		c:    c,
		next: next,
		name: name,
	}, nil
}

const (
	contentLength    = "Content-Length"
	contentType      = "Content-Type"
	cacheStatus      = "Cache-Status"
	cacheHitStatus   = "hit"
	cacheMissStatus  = "miss"
	cacheErrorStatus = "error"
)

func (a *ImageOptimizer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// TODO Check if cacheable

	targetFormat := "image/webp"
	// Return cached result here.
	key, err := cache.Tokenize(req)
	if err != nil {
		panic(err)
	}

	if v, err := a.c.Get(key); err == nil {
		rw.Header().Set(cacheStatus, cacheHitStatus)
		rw.Header().Set(contentLength, fmt.Sprint(len(v)))
		rw.Header().Add(contentType, targetFormat)
		_, err = rw.Write(v)
		if err != nil {
			panic(err)
		}
		return
	}

	wrappedWriter := &responseWriter{
		ResponseWriter: rw,
		bypassHeader:   true,
	}

	a.next.ServeHTTP(wrappedWriter, req)

	wrappedWriter.bypassHeader = false
	bodyBytes := wrappedWriter.buffer.Bytes()

	// If not image response, forward original and leave it here.
	if !isImageResponse(rw) {
		_, err = rw.Write(bodyBytes)
		if err != nil {
			panic(err)
		}
		return
	}

	width, err := imageWidthRequest(req)
	if err != nil {
		panic(err)
	}

	optimized, ct, err := a.p.Optimize(bodyBytes, rw.Header().Get(contentType), targetFormat, 75, width)
	if err != nil {
		panic(err)
	}

	rw.Header().Set(contentLength, fmt.Sprint(len(optimized)))
	rw.Header().Set(contentType, ct)
	rw.Header().Set(cacheStatus, cacheMissStatus)

	_, err = rw.Write(optimized)
	if err != nil {
		panic(err)
	}

	err = a.c.Set(key, optimized, 100*time.Second)
	if err != nil {
		panic(err)
	}
}

func imageWidthRequest(req *http.Request) (int, error) {
	w := req.URL.Query().Get("w")

	// if no query param
	if len(w) == 0 {
		return 0, nil
	}

	v, err := strconv.Atoi(w)
	if err != nil {
		return 0, err
	}

	if v < 0 {
		return 0, errors.New("width cannot be negative value")
	}

	return v, nil
}

// isImageResponse Determine with Content-Type header if the response is an image.
func isImageResponse(rw http.ResponseWriter) bool {
	return strings.HasPrefix(rw.Header().Get(contentType), "image/")
}

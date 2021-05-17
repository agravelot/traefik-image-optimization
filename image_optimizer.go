// Package imageopti Image optimizer plugin
package imageopti

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/agravelot/imageopti/cache"
	"github.com/agravelot/imageopti/config"
	"github.com/agravelot/imageopti/processor"
)

// Config the plugin configuration.
type Config struct {
	config.Config
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		config.Config{
			Processor: "",
			Cache:     "",
			Imaginary: config.ImaginaryProcessorConfig{URL: ""},
			Redis:     config.RedisCacheConfig{URL: ""},
			File:      config.FileCacheConfig{Path: ""},
			Picture:   config.PictureProcessingConfig{Formats: make(map[string]config.PictureFormat)},
		},
	}
}

// ImageOptimizer middleware plugin struct.
type ImageOptimizer struct {
	next http.Handler
	name string
	p    processor.Processor
	c    cache.Cache

	formatRegExp        *regexp.Regexp
	formatRegExpReplace string
	formats             map[string]config.PictureFormat
}

// New created a new ImageOptimizer plugin.
func New(ctx context.Context, next http.Handler, conf *Config, name string) (http.Handler, error) {
	log.Println("Loading image optimization plugin...")

	if conf.Processor == "" {
		return nil, fmt.Errorf("processor must be defined")
	}

	c, err := cache.New(conf.Config)
	if err != nil {
		panic(err)
	}

	p, err := processor.New(conf.Config)
	if err != nil {
		panic(err)
	}

	return &ImageOptimizer{
		p:    p,
		c:    c,
		next: next,
		name: name,

		formatRegExp:        getFormatRegExp(conf.Picture.FormatRegExp),
		formatRegExpReplace: conf.Picture.FormatRegExpReplace,
		formats:             conf.Picture.Formats,
	}, nil
}

const (
	contentLength   = "Content-Length"
	contentType     = "Content-Type"
	cacheStatus     = "Cache-Status"
	cacheHitStatus  = "hit"
	cacheMissStatus = "miss"
	cacheExpiry     = 100 * time.Second
	targetFormat    = "image/webp"
)

func (a *ImageOptimizer) transformPath(req *http.Request) (width int) {
	log.Println(req.URL.Path, a.formatRegExp.MatchString(req.URL.Path))
	if a.formatRegExp == nil || !a.formatRegExp.MatchString(req.URL.Path) {
		return 0
	}
	path := string(a.formatRegExp.ReplaceAllFunc([]byte(req.URL.Path), func(s []byte) []byte {
		format := a.formatRegExp.ReplaceAllString(string(s), "$1")
		width, _ = strconv.Atoi(format)
		if width == 0 {
			formatConfig := a.formats[format]
			width = formatConfig.Width
		}
		return []byte(a.formatRegExpReplace)
	}))

	var err error
	req.URL.Path, err = url.PathUnescape(path)
	if err != nil {
		panic(err)
	}
	req.RequestURI = req.URL.RequestURI()
	return
}

func (a *ImageOptimizer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// TODO Check if cacheable
	key, err := cache.Tokenize(req)
	if err != nil {
		panic(err)
	}
	// Return cached result here.
	if v, err := a.c.Get(key); err == nil {
		rw.Header().Set(contentLength, fmt.Sprint(len(v)))
		rw.Header().Set(contentType, targetFormat)
		rw.Header().Set(cacheStatus, cacheHitStatus)
		_, err = rw.Write(v)

		if err != nil {
			panic(err)
		}

		return
	}

	wrappedWriter := &responseWriter{
		ResponseWriter: rw,
		bypassHeader:   true,
		wroteHeader:    false,
		buffer:         bytes.Buffer{},
	}

	width := a.transformPath(req)
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

	if width == 0 {
		width, err = imageWidthRequest(req)
		if err != nil {
			panic(err)
		}
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

	err = a.c.Set(key, optimized, cacheExpiry)
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
		return 0, fmt.Errorf("unable to convert w query param to int: %w", err)
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

func getFormatRegExp(rx string) *regexp.Regexp {
	if rx != "" {

		rx, err := regexp.Compile(rx)
		if err != nil {
			panic(err)
		}
		return rx
	}
	return nil
}

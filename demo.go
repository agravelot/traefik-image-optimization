// Package plugindemo a demo plugin.
package plugindemo

import (
	"context"
	"log"
	"net/http"
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
	// Add logic here

	a.next.ServeHTTP(rw, req)
}

package cache

import (
	"fmt"
	"net/http"
)

// Tokenize generate unique key for request caching strategy.
func Tokenize(req *http.Request) (string, error) {
	width := req.URL.Query().Get("w")

	if len(width) == 0 {
		width = "original"
	}

	return fmt.Sprintf("%s:%s:%s:%s:%s", req.Method, req.URL.Scheme, req.Host, req.URL.Path, width), nil
}

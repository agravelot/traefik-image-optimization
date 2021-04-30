package plugindemo_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/agravelot/plugindemo"
)

func TestDemo(t *testing.T) {
	cfg := plugindemo.CreateConfig()

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := plugindemo.New(ctx, next, cfg, "demo-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	// TODO Assert here
}

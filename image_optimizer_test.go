package image_optimizer

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/agravelot/image_optimizer/config"
)

func TestImageOptimizer_ServeHTTP(t *testing.T) {
	type args struct {
		config config.Config
	}
	tests := []struct {
		name                      string
		args                      args
		remoteResponseContentType string
		remoteResponseContent     []byte
		want                      bool
		wantErr                   bool
	}{
		{
			name:                      "should pass with processor",
			args:                      args{config: config.Config{Processor: "imaginary", Imaginary: config.ImaginaryProcessorConfig{Url: "http://localhost"}}},
			want:                      false,
			wantErr:                   false,
			remoteResponseContentType: "text/html",
			remoteResponseContent:     []byte("dummy response"),
		},
		{
			name:                      "should not pass without processor",
			args:                      args{config: config.Config{Processor: ""}},
			want:                      false,
			wantErr:                   true,
			remoteResponseContentType: "text/html",
			remoteResponseContent:     []byte("dummy response"),
		},
		{
			name:                      "should pass with processor with image response",
			args:                      args{config: config.Config{Processor: "none"}},
			want:                      false,
			wantErr:                   false,
			remoteResponseContentType: "image/jpeg",
			remoteResponseContent:     []byte("dummy image"),
		},
		{
			name:                      "should return original response if not image request",
			args:                      args{config: config.Config{Processor: "none"}},
			want:                      false,
			wantErr:                   false,
			remoteResponseContentType: "text/html",
			remoteResponseContent:     []byte("dummy response"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := CreateConfig()
			cfg.Processor = tt.args.config.Processor
			cfg.Imaginary = tt.args.config.Imaginary

			ctx := context.Background()
			next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				rw.Header().Add("content-type", tt.remoteResponseContentType)
				rw.Write([]byte(tt.remoteResponseContent))
			})

			handler, err := New(ctx, next, cfg, "demo-plugin")

			if (err != nil) != tt.wantErr {
				t.Fatalf("New() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil {
				return
			}

			recorder := httptest.NewRecorder()

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
			if err != nil {
				t.Fatal(err)
			}

			handler.ServeHTTP(recorder, req)

			if !bytes.Equal(recorder.Body.Bytes(), tt.remoteResponseContent) {
				t.Fatal("response are not equals")
			}
		})
	}
}

func TestIsImageResponse(t *testing.T) {
	type args struct {
		contentType string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "should return false with empty string",
			args:    args{contentType: ""},
			want:    false,
			wantErr: false,
		},
		{
			name:    "should return true with jpeg content type",
			args:    args{contentType: "image/jpeg"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "should return true with webp content type",
			args:    args{contentType: "image/webp"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "should return false with json content type",
			args:    args{contentType: "application/json"},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			recorder.Header().Add("Content-Type", tt.args.contentType)

			got := isImageResponse(recorder)

			if got != tt.want {
				t.Errorf("IsImageResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImageWidthRequest(t *testing.T) {
	type args struct {
		url string
	}
	ctx := context.Background()

	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "should return error with positve width",
			args:    args{url: "http://localhost?w=124"},
			want:    124,
			wantErr: false,
		},
		{
			name:    "should return error with negative width",
			args:    args{url: "http://localhost?w=-124"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "should return error with text as width",
			args:    args{url: "http://localhost?w=azeaze"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "should return error with empty width",
			args:    args{url: "http://localhost?w="},
			want:    0,
			wantErr: false,
		},
		{
			name:    "should return error with no width",
			args:    args{url: "http://localhost"},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, tt.args.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			got, err := imageWidthRequest(req)

			if (err != nil) != tt.wantErr {
				t.Fatalf("imageWidthRequest() error = %v, wantErr %v", err, tt.wantErr)
			}

			if got != tt.want {
				t.Errorf("ImageWidthRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

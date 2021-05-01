package image_optimizer_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/agravelot/image_optimizer"
)

func TestDemo(t *testing.T) {
	type args struct {
		config image_optimizer.Config
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "should init with processor",
			args:    args{config: image_optimizer.Config{Processor: "imaginary"}},
			want:    false,
			wantErr: false,
		},
		{
			name:    "should not init without processor",
			args:    args{config: image_optimizer.Config{Processor: ""}},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := image_optimizer.CreateConfig()
			cfg.Processor = tt.args.config.Processor

			ctx := context.Background()
			next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

			handler, err := image_optimizer.New(ctx, next, cfg, "demo-plugin")

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

			// TODO Assert here
		})
	}
}

func TestIsImageRequest(t *testing.T) {
	type args struct {
		acceptHeader string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "should return false with empty string",
			args:    args{acceptHeader: ""},
			want:    false,
			wantErr: false,
		},
		{
			name:    "should return false with standard chrome request",
			args:    args{acceptHeader: "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
			want:    false,
			wantErr: false,
		},
		{
			name:    "should return true with image chrome request",
			args:    args{acceptHeader: "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8"},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := image_optimizer.IsImageRequest(tt.args.acceptHeader)

			if got != tt.want {
				t.Errorf("IsImageRequest() = %v, want %v", got, tt.want)
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
			got := image_optimizer.IsImageResponse(tt.args.contentType)

			if got != tt.want {
				t.Errorf("IsImageResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

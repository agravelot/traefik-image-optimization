package cache

import (
	"context"
	"net/http"
	"testing"
)

func TestTokenize(t *testing.T) {
	ctx := context.Background()

	newRequest := func(method, url string) *http.Request {
		req, err := http.NewRequestWithContext(ctx, method, url, nil)
		if err != nil {
			t.Fatal(err)
		}
		return req
	}

	type args struct {
		req *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "should return correct token",
			args:    args{req: newRequest(http.MethodGet, "http://localhost/img.jpeg")},
			want:    "GET:http:localhost:/img.jpeg:original",
			wantErr: false,
		},
		{
			name:    "should return correct token with width query param",
			args:    args{req: newRequest(http.MethodGet, "http://localhost/img.jpeg?w=1024")},
			want:    "GET:http:localhost:/img.jpeg:1024",
			wantErr: false,
		},
		{
			name:    "should return correct token with width query param",
			args:    args{req: newRequest(http.MethodDelete, "http://localhost/img.jpeg?w=1024")},
			want:    "DELETE:http:localhost:/img.jpeg:1024",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Tokenize(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tokenize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Tokenize() = %v, want %v", got, tt.want)
			}
		})
	}
}

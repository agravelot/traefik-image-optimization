package processor

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/agravelot/imageopti/config"
)

const defaultTimeout = 1 * time.Second

func TestNew(t *testing.T) {
	type args struct {
		conf config.Config
	}

	tests := []struct {
		name    string
		args    args
		want    Processor
		wantErr bool
	}{
		{
			name: "should be able to return imaginary optimizer",
			args: args{
				config.Config{
					Processor: "imaginary",
					Imaginary: config.ImaginaryProcessorConfig{URL: "http://localhost"},
					Cache:     "none",
					Redis:     config.RedisCacheConfig{URL: ""},
					File:      config.FileCacheConfig{Path: ""},
				},
			},
			want:    &ImaginaryProcessor{"", http.Client{Timeout: defaultTimeout}},
			wantErr: false,
		},
		{
			name:    "should not be able to init imaginary without valid url",
			args:    args{config.Config{Processor: "imaginary", Imaginary: config.ImaginaryProcessorConfig{URL: "localhost"}}},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should not be able to init imaginary without valid url 2 ",
			args: args{
				config.Config{
					Processor: "imaginary",
					Imaginary: config.ImaginaryProcessorConfig{URL: "htt://localhost"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "should not be able to init imaginary without url",
			args:    args{config.Config{Processor: "imaginary"}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "should be able to return local optimizer",
			args:    args{config.Config{Processor: "local"}},
			want:    &LocalProcessor{},
			wantErr: false,
		},
		{
			name:    "should return error with unsupported processor",
			args:    args{config.Config{Processor: "unsupported"}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "should return error with empty processor",
			args:    args{config.Config{Processor: "unsupported"}},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Fatalf("New() error = %v, wantErr %v", err, tt.wantErr)
			}

			typeGot := reflect.TypeOf(got)
			typeWanted := reflect.TypeOf(tt.want)

			if typeGot != typeWanted {
				t.Errorf("New() = %v, want %v", typeGot, typeWanted)
			}
		})
	}
}

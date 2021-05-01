package optimizer_test

import (
	"reflect"
	"testing"

	"github.com/agravelot/plugindemo/optimizer"
)

func TestNew(t *testing.T) {
	type args struct {
		driver string
	}
	tests := []struct {
		name    string
		args    args
		want    optimizer.Optimizer
		wantErr bool
	}{
		{
			name:    "should be able to return imaginary optimizer",
			args:    args{driver: "imaginary"},
			want:    &optimizer.ImaginaryOptimizer{},
			wantErr: false,
		},
		{
			name:    "should be able to return local optimizer",
			args:    args{driver: "local"},
			want:    &optimizer.LocalOptimizer{},
			wantErr: false,
		},
		{
			name:    "should return error with unsupported driver",
			args:    args{driver: "unsupported"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := optimizer.New(tt.args.driver)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			typeGot := reflect.TypeOf(got)
			typeWanted := reflect.TypeOf(tt.want)

			if typeGot != typeWanted {
				t.Errorf("New() = %v, want %v", typeGot, typeWanted)
			}
		})
	}
}

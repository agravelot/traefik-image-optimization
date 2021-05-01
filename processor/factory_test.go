package processor_test

import (
	"reflect"
	"testing"

	"github.com/agravelot/image_optimizer/processor"
)

func TestNew(t *testing.T) {
	type args struct {
		driver string
	}
	tests := []struct {
		name    string
		args    args
		want    processor.Optimizer
		wantErr bool
	}{
		{
			name:    "should be able to return imaginary optimizer",
			args:    args{driver: "imaginary"},
			want:    &processor.ImaginaryOptimizer{},
			wantErr: false,
		},
		{
			name:    "should be able to return local optimizer",
			args:    args{driver: "local"},
			want:    &processor.LocalOptimizer{},
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
			got, err := processor.New(tt.args.driver)
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

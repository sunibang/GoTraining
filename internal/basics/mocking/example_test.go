package mocking

import (
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/romangurevitch/go-training/internal/basics/mocking/calculator"
	"github.com/romangurevitch/go-training/internal/basics/mocking/calculator/mocks"
)

func TestExampleFunction(t *testing.T) {
	type args struct {
		adder func() calculator.Adder
		x     int
		y     int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "no error testcase",
			args: args{
				adder: func() calculator.Adder {
					ctrl := gomock.NewController(t)
					a := mocks.NewMockAdder(ctrl)
					a.EXPECT().SingleDigitAdd(1, 2).Return(3, nil)
					return a
				},
				x: 1,
				y: 2,
			},
			want:    3,
			wantErr: false,
		},
		{
			name: "error testcase",
			args: args{
				adder: func() calculator.Adder {
					ctrl := gomock.NewController(t)
					a := mocks.NewMockAdder(ctrl)
					a.EXPECT().SingleDigitAdd(gomock.Any(), gomock.Any()).Return(0, errors.New("error"))
					return a
				},
				x: 1,
				y: 2,
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExampleFunction(tt.args.adder(), tt.args.x, tt.args.y)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExampleFunction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExampleFunction() got = %v, want %v", got, tt.want)
			}
		})
	}
}

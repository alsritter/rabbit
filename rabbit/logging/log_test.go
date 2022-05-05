package logging

import "testing"

func TestError(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"test error",
			args{
				v: "this is a test error msg ~",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Error(tt.args.v)
		})
	}
}

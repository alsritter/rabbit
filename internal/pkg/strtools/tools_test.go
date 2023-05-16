package strtools

import "testing"

func TestLimit(t *testing.T) {
	type args struct {
		str    string
		joint  string
		length int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				str:    "这是一个君伦的上门服务",
				joint:  "...",
				length: 3,
			},
			want: "这是一...",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Limit(tt.args.str, tt.args.joint, tt.args.length); got != tt.want {
				t.Errorf("Limit() = %v, want %v", got, tt.want)
			}
		})
	}
}

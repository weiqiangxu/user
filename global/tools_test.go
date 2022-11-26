package global

import (
	"testing"
)

func TestGenerateUniqueId(t *testing.T) {
	type args struct {
		size int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test 16bit unique id",
			args: args{
				size: 16,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateUniqueId(tt.args.size)
			t.Logf("%s", got)
		})
	}
}

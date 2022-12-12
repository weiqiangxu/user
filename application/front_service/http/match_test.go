package http

import "testing"

func Test_validateStruct(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "test",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateStruct()
			t.Logf("%s", err.Error())
		})
	}
}

func Test_validateStruct2(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "验证自定义错误",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateStruct2()
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

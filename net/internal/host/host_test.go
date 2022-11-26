package host

import "testing"

func TestExtractHostPort(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name     string
		args     args
		wantHost string
		wantPort uint64
		wantErr  bool
	}{
		{
			name: "test extract host port",
			args: args{
				addr: "www.baidu.com",
			},
			wantHost: "",
			wantPort: 0,
			wantErr:  false,
		},
		{
			name: "test extract host port",
			args: args{
				addr: "192.168.1.1:8989",
			},
			wantHost: "",
			wantPort: 0,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHost, gotPort, err := ExtractHostPort(tt.args.addr)
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("host=%#v", gotHost)
			t.Logf("port=%#v", gotPort)
		})
	}
}

func Test_isValidIP(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test ip",
			args: args{
				addr: "192.168.1.1",
			},
			want: false,
		},
		{
			name: "test ip",
			args: args{
				addr: "0.999999.0.1",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := isValidIP(tt.args.addr)
			t.Logf("valid=%#v", valid)
		})
	}
}

package smtp

import "testing"

func Test_extractServiceId(t *testing.T) {
	type args struct {
		to string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "regular",
			args: args{
				to: "telegram@mydomain.tld",
			},
			want: "telegram",
		},
		{
			name: "local delivery only",
			args: args{
				to: "telegram",
			},
			want: "telegram",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractServiceId(tt.args.to); got != tt.want {
				t.Errorf("extractServiceId() = %v, want %v", got, tt.want)
			}
		})
	}
}

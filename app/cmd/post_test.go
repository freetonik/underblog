package cmd

import "testing"

func Test_fNameWithoutExtension(t *testing.T) {
	type args struct {
		fn string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"24-10-2019-Welcome.md", args{fn: "24-10-2019-Welcome.md"}, "24-10-2019-Welcome"},
		{"24-10-2019-Welcome", args{fn: "24-10-2019-Welcome"}, "24-10-2019-Welcome"},
		{"24-10-2019-Welcome.wtf.md", args{fn: "24-10-2019-Welcome.wtf.md"}, "24-10-2019-Welcome.wtf"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fNameWithoutExtension(tt.args.fn); got != tt.want {
				t.Errorf("fNameWithoutExtension() = %v, want %v", got, tt.want)
			}
		})
	}
}

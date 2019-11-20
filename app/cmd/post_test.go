package cmd

import (
	. "github.com/onsi/gomega"
	"gopkg.in/russross/blackfriday.v2"
	"testing"
)

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

func TestHighlightCode(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{
			text: "```ruby\nputs 'hello world'\n```\n",
			want: "<pre><code class=\"language-ruby\"><span class=\"pln\">puts</span> <span class=\"str\">&#39;hello world&#39;</span>\n</code></pre>\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGomegaWithT(t)
			markdown := blackfriday.Run([]byte(tt.text))
			result, err := HighlightCode(markdown)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(string(result)).To(Equal(tt.want))
		})
	}
}

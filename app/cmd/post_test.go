package cmd

import (
	"errors"
	. "github.com/onsi/gomega"
	"testing"
	"time"
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

func TestExtractMetaFromFilename(t *testing.T) {
	tests := []struct {
		filename string
		want     Post
		err      error
	}{
		{
			"2019-11-17-Welcome.md",
			Post{"", "", time.Date(2019, 11, 17, 0, 0, 0, 0, time.UTC), "Welcome"},
			nil,
		},
		{
			"FilenameWithoutDate.md",
			Post{},
			errors.New("can't parse filename 'FilenameWithoutDate.md', it should be in format 'YYYY-MM-DD-slug.md'"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			g := NewGomegaWithT(t)

			post, err := ExtractMetaFromFilename(tt.filename)

			if tt.err != nil {
				g.Expect(err).To(Equal(tt.err))
				g.Expect(post).To(Equal(Post{}))
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(post.Date).To(Equal(tt.want.Date))
				g.Expect(post.Slug).To(Equal(tt.want.Slug))
			}
		})
	}
}

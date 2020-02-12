package cmd

import (
	"testing"
)

func TestNewRSS(t *testing.T) {
	blog := &Blog{}
	rss := NewRSS(blog)
	if rss.blog != blog {
		t.Errorf("newRss() fails. Expected %p, got %p", blog, rss.blog)
	}
}

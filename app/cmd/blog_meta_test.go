package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlogLink(t *testing.T) {
	meta := BlogMeta{}
	res := meta.BlogLink("test")
	assert.Equal(t, "test", res)
	assert.Equal(t, "test", meta.Link)
}

func TestBlogTitle(t *testing.T) {
	meta := BlogMeta{}
	res := meta.BlogTitle("test")
	assert.Equal(t, "test", res)
	assert.Equal(t, "test", meta.Title)
}

func TestBlogDescription(t *testing.T) {
	meta := BlogMeta{}
	res := meta.BlogDescription("test")
	assert.Equal(t, "test", res)
	assert.Equal(t, "test", meta.Description)
}

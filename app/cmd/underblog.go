package cmd

import (
	"fmt"

	"github.com/freetonik/underblog/app/internal"
)

// MakeBlog initiates blog rendering process
func MakeBlog(opts internal.Opts) error {
	blog := NewBlog(opts)

	err := blog.Render()
	if err != nil {
		return fmt.Errorf("can't render html: %v", err)
	}

	// support old template version without blog URL
	if blog.meta.Link != "" {
		if err := NewRSS(blog).Render(); err != nil {
			return fmt.Errorf("can't render RSS: %v", err)
		}
	}

	return nil
}

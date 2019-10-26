package cmd

import (
	"errors"
	"fmt"
	"github.com/freetonik/underblog/app/internal"
)

func MakeBlog(opts internal.Opts) error {
	blog := NewBlog(opts)

	err := blog.Render()
	if err != nil {
		return errors.New(fmt.Sprintf("can't render html: %v", err))
	}

	return nil
}

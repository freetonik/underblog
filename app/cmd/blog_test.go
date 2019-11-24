package cmd_test

import (
	. "github.com/freetonik/underblog/app/cmd"
	"github.com/freetonik/underblog/app/internal"
	"reflect"
	"testing"
	"time"
)

func TestBlog_SortPosts(t *testing.T) {
	blog := NewBlog(internal.Opts{})

	date1 := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)   // Jan 1
	date2 := time.Date(2019, 1, 3, 0, 0, 0, 0, time.UTC)   // Jan 3
	date3 := time.Date(2019, 1, 2, 0, 0, 0, 0, time.UTC)   // Jan 2
	date4 := time.Date(2018, 12, 31, 0, 0, 0, 0, time.UTC) // Dec 31

	post1 := Post{Date: date1}
	post2 := Post{Date: date2}
	post3 := Post{Date: date3}
	post4 := Post{Date: date4}

	blog.Posts = append(blog.Posts, post1, post2, post3, post4)

	blog.SortPosts()

	expectedOrder := []time.Time{date2, date3, date1, date4}
	resultedOrder := []time.Time{blog.Posts[0].Date, blog.Posts[1].Date, blog.Posts[2].Date, blog.Posts[3].Date}

	if !reflect.DeepEqual(resultedOrder, expectedOrder) {
		t.Errorf("Blog.SortPost() fails. Expected %s, got %s", expectedOrder, resultedOrder)
	}
}

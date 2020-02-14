package cmd

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func testGetStubMeta() BlogMeta {
	return BlogMeta{
		Link:        "prefix",
		Title:       "title",
		Description: "description",
	}
}

func testGetStubPosts() []Post {
	return []Post{
		Post{Title: "#1", Body: template.HTML("Post body #1"), Date: time.Now(), Slug: "post-1"},
		Post{Title: "#2", Body: template.HTML("Post body #2"), Date: time.Now(), Slug: "post-2"},
		Post{Title: "#3", Body: template.HTML("Post body #3"), Date: time.Now(), Slug: "post-3"},
	}
}

func TestNewRSS(t *testing.T) {
	blog := &Blog{}
	rss := NewRSS(blog)
	assert.Same(t, blog, rss.blog)
}

func TestRender(t *testing.T) {
	blog := &Blog{Posts: testGetStubPosts(), meta: testGetStubMeta()}
	rss := NewRSS(blog)

	err := rss.Render("")
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "open : no such file or directory", err.Error())
	}

	tmp, _ := ioutil.TempFile("/tmp", "")
	path := tmp.Name()
	defer os.Remove(path)

	err = rss.Render(path)
	assert.Nil(t, err)
	assert.FileExists(t, path)

	stat, _ := os.Stat(path)
	assert.NotEmpty(t, stat.Size)
}

func TestGetFeedXML(t *testing.T) {
	post := testGetStubPosts()[0]
	blog := &Blog{Posts: []Post{post}, meta: testGetStubMeta()}
	rss := NewRSS(blog)

	expected := fmt.Sprintf(`
		<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
			<channel>
				<title>title</title>
				<link>prefix</link>
				<description>description</description>
				<atom:link href="prefix/rss.xml" rel="self" type="application/rss+xml" />
				<item>
					<title>#1</title>
					<link>prefix/posts/post-1</link>
					<guid>prefix/posts/post-1</guid>
					<description>Post body #1</description>
					<pubDate>%s</pubDate>
				</item>
			</channel>
    </rss>
	`, post.Date.Format(time.RFC1123Z))

	res := rss.getFeedXML()
	assert.NotEmpty(t, res)
	assert.Equal(
		t,
		strings.Join(strings.Fields(expected), ""),
		strings.Join(strings.Fields(string(res)), ""),
	)
}

func TestGetFeed(t *testing.T) {
	blog := &Blog{Posts: testGetStubPosts(), meta: testGetStubMeta()}
	rss := NewRSS(blog)

	res := rss.getFeed()
	assert.Equal(t, "2.0", res.Version)
	assert.Equal(t, "http://www.w3.org/2005/Atom", res.ContentNamespace)
}

func TestGetFeedChannel(t *testing.T) {
	blog := &Blog{Posts: testGetStubPosts(), meta: testGetStubMeta()}
	rss := NewRSS(blog)

	res := rss.getFeedChannel()
	assert.Equal(t, blog.meta.Title, res.Title)
	assert.Equal(t, blog.meta.Description, res.Description)
	assert.Equal(t, blog.meta.Link, res.Link)
	assert.Equal(t, blog.meta.Link+"/rss.xml", res.AtomLink.Link)
	assert.Equal(t, "self", res.AtomLink.Rel)
	assert.Equal(t, "application/rss+xml", res.AtomLink.Type)
	assert.Len(t, res.Items, len(blog.Posts))
}

func TestGetFeedItems(t *testing.T) {
	posts := testGetStubPosts()
	blog := &Blog{Posts: posts, meta: testGetStubMeta()}
	rss := NewRSS(blog)

	res := rss.getFeedItems()
	assert.Len(t, res, len(posts))

	for i, item := range res {
		post := posts[i]

		url := fmt.Sprintf("%s/posts/%s", blog.meta.Link, post.Slug)
		assert.Equal(t, url, item.Link)
		assert.Equal(t, post.Title, item.Title)
		assert.Equal(t, post.Date.Format(time.RFC1123Z), item.Published)
	}
}

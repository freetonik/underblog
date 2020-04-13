package cmd

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
	"time"
)

// RSS represents blog feed structure
type RSS struct {
	blog *Blog
}

type feed struct {
	XMLName          xml.Name `xml:"rss"`
	Version          string   `xml:"version,attr"`
	ContentNamespace string   `xml:"xmlns:atom,attr"`
	Channel          *feedChannel
}

type feedAtomLink struct {
	XMLName xml.Name `xml:"atom:link,allowempty"`
	Link    string   `xml:"href,attr"`
	Rel     string   `xml:"rel,attr"`
	Type    string   `xml:"type,attr"`
}

type feedChannel struct {
	XMLName     xml.Name `xml:"channel"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	AtomLink    *feedAtomLink
	Items       []*feedItem `xml:"item"`
}

type feedItem struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	GUID        string   `xml:"guid"`
	Description string   `xml:"description"`
	Published   string   `xml:"pubDate"`
}

// NewRSS initialize new RSS structure
func NewRSS(blog *Blog) *RSS {
	rss := new(RSS)
	rss.blog = blog

	return rss
}

// Render generates blog RSS feed and store it in a file
func (r *RSS) Render(path string) error {
	feed := r.getFeedXML()

	return ioutil.WriteFile(path, feed, os.ModePerm)
}

func (r *RSS) getFeedXML() []byte {
	feed := r.getFeed()

	// skip an error check 'cause it will always return some result
	out, _ := xml.MarshalIndent(feed, "  ", "    ")

	// "dirty" hack for get self-closing tag
	res := strings.Replace(string(out), "></atom:link>", " />", 1)
	return []byte(res)
}

func (r *RSS) getFeed() *feed {
	return &feed{
		Version:          "2.0",
		Channel:          r.getFeedChannel(),
		ContentNamespace: "http://www.w3.org/2005/Atom",
	}
}

func (r *RSS) getFeedChannel() *feedChannel {
	meta := r.blog.meta
	return &feedChannel{
		Title:       meta.Title,
		Description: meta.Description,
		Link:        meta.Link,
		Items:       r.getFeedItems(),
		AtomLink: &feedAtomLink{
			Link: r.blog.meta.Link + "/rss.xml",
			Rel:  "self",
			Type: "application/rss+xml",
		},
	}
}

func (r *RSS) getFeedItems() []*feedItem {
	items := make([]*feedItem, len(r.blog.Posts))
	for i, post := range r.blog.Posts {
		link := r.blog.meta.Link + post.getURL()
		items[i] = &feedItem{
			Title:       post.Title,
			Description: template.HTMLEscapeString(string(post.Body)),
			Link:        link,
			GUID:        link,
			Published:   post.Date.Format(time.RFC1123Z),
		}
	}
	return items
}

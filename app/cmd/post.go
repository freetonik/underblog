package cmd

import (
	"fmt"
	"gopkg.in/russross/blackfriday.v2"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func NewPost(filename string) Post {
	// Get filename without extension
	filenameBase := fNameWithoutExtension(filename)
	verifyFilenameBaseFormat(filenameBase)

	// Get date and slug from filename
	day := filenameBase[0:2]
	month := filenameBase[3:5]
	year := filenameBase[6:10]
	date, err := time.Parse("02-01-2006", day+"-"+month+"-"+year)
	if err != nil {
		log.Fatal(err)
	}
	slug := filenameBase[11:]

	// Get body from file
	mdfile, err := os.Open("./markdown/" + filename)
	if err != nil {
		log.Fatal(err)
	}
	rawBytes, err := ioutil.ReadAll(mdfile)

	// Get title from first line of file
	lines := strings.Split(string(rawBytes), "\n")
	title := strings.Replace(lines[0], "# ", "", -1)

	_ = mdfile.Close()

	// Convert Markdown to HTML
	html := blackfriday.Run(rawBytes)

	// Create a Post struct
	post := Post{
		Title: title,
		Body:  template.HTML(html),
		Date:  date,
		Slug:  slug,
	}

	// Save file
	post.createFile()

	return post
}

type Post struct {
	Title string
	Body  template.HTML
	Date  time.Time
	Slug  string
}

func (post *Post) createFile() {
	// Create folder for HTML
	newPath := filepath.Join("public/posts", post.Slug)
	_ = os.MkdirAll(newPath, os.ModePerm)

	// Create HTML file
	f, err := os.Create("public/posts/" + post.Slug + "/" + "index.html")
	if err != nil {
		log.Fatal(err)
	}

	// Generate final HTML file from template
	t, _ := template.ParseFiles("post.html")
	err = t.Execute(f, post)
	if err != nil {
		log.Fatalf("can't execute template: %v", err)
	}
	_ = f.Close()
}

func fNameWithoutExtension(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}

func verifyFilenameBaseFormat(f string) {
	errorDescription := "I can't parse this filename. Make sure its name is formatted as: DD-MM-YYY-slug.md"

	if len(f) < 12 {
		fmt.Println(errorDescription)
		os.Exit(1)
	}

	// day is int?
	_, err := strconv.Atoi(f[0:2])
	if err != nil {
		fmt.Println(errorDescription)
		os.Exit(1)
	}

	// month is int?
	_, err2 := strconv.Atoi(f[3:5])
	if err2 != nil {
		fmt.Println(errorDescription)
		os.Exit(1)
	}

	// year is int?
	_, err3 := strconv.Atoi(f[6:10])
	if err3 != nil {
		fmt.Println(errorDescription)
		os.Exit(1)
	}
}

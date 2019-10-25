package main

import (
	"flag"
	"fmt"
	"gopkg.in/russross/blackfriday.v2"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Post struct {
	Title string
	Body  template.HTML
	Date  time.Time
	Slug  string
}

func init() {
	const AppVersion = "0.1.2"

	version := flag.Bool("version", false, "prints current roxy version")
	flag.Parse()
	if *version {
		fmt.Println(AppVersion)
		os.Exit(0)
	}
}

func isFileValid(file os.FileInfo) bool {
	return path.Ext(file.Name()) == ".md" || path.Ext(file.Name()) == ".markdown"
}

func main() {
	start := time.Now()

	fmt.Printf("Starting...\n")

	var posts []Post

	// Read markdown files folder
	files, err := ioutil.ReadDir("./markdown/")
	if err != nil {
		fmt.Println("I need a folder named 'markdown' to continue")
		log.Fatal(err)
	}

	// For each file, create HTML
	wg := &sync.WaitGroup{}
	for _, file := range files {
		if isFileValid(file) {
			fmt.Println("Processing " + file.Name())

			wg.Add(1)
			post := createPost(file.Name())
			go createPostFile(post, wg)

			posts = append(posts, post)

			fmt.Println("Done with  " + file.Name())
			fmt.Println("---")
		}
	}
	wg.Wait()

	// Create blog root HTML
	newPath := filepath.Join(".", "public")
	_ = os.MkdirAll(newPath, os.ModePerm)

	f, err := os.Create("public/index.html")
	if err != nil {
		log.Fatal(err)
	}

	t, _ := template.ParseFiles("index.html")
	err = t.Execute(f, posts)
	if err != nil {
		log.Fatalf("can't execute template: %v", err)
	}
	_ = f.Close()

	// Copy styles
	copyCssToPublicDir()

	elapsed := time.Since(start)
	log.Printf("Done in %s", elapsed)
}

func createPost(filename string) Post {
	// Get filename without extension
	filenameBase := FnameWithoutExtension(filename)
	VerifyFilenameBaseFormat(filenameBase)

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

	return post
}

func createPostFile(post Post, wg *sync.WaitGroup) {
	defer wg.Done()
	// Create folder for HTML
	newPath := filepath.Join("public/posts", post.Slug)
	_ = os.MkdirAll(newPath, os.ModePerm)

	// Create HTML file
	f, err := os.Create("public/posts/" + post.Slug + "/" + "index.html")
	if err != nil {
		fmt.Println("Aaa!")
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

func FnameWithoutExtension(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}

func VerifyFilenameBaseFormat(f string) {
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

func copyCssToPublicDir() {
	from, err := os.Open("./css/styles.css")
	if err != nil {
		log.Fatal(err)
	}

	newPath := filepath.Join("public", "css")
	_ = os.MkdirAll(newPath, os.ModePerm)

	to, err := os.OpenFile("./public/css/styles.css", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(to, from)
	if err != nil {
		log.Fatal(err)
	}

	_ = from.Close()
	_ = to.Close()
}

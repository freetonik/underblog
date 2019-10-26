package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"gopkg.in/russross/blackfriday.v2"
)

var dateLayout = "2006-01-02"

type Post struct {
	Title string
	Body  template.HTML
	Date  time.Time
	Slug  string
}

func init() {
	const AppVersion = "0.1.2"

	version := flag.Bool("version", false, "prints current roxy version")
	flag.StringVar(&dateLayout, "date-layout", dateLayout, "date layout for file names")
	flag.Parse()
	if *version {
		log.Print(AppVersion)
		os.Exit(0)
	}
}

func main() {
	start := time.Now()
	log.Print("Starting...")

	var posts []*Post

	// Read markdown files folder
	files, err := ioutil.ReadDir("./markdown/")
	if err != nil {
		log.Fatalf("I need a folder named 'markdown' to continue: %s", err)
	}

	wg := &sync.WaitGroup{}

	// For each file, create HTML
	for _, file := range files {
		if path.Ext(file.Name()) == ".md" || path.Ext(file.Name()) == ".markdown" {
			log.Printf("Processing %s", file.Name())
			post, err := createPost(file.Name())
			if err != nil {
				log.Fatalf("Creating Post failed: %s", err)
			}
			wg.Add(1)
			go func() {
				if err := createPostFile(post); err != nil {
					log.Fatalf("Creating Post file failed: %s", err)
				}
				wg.Done()
			}()
			posts = append(posts, post)
			log.Println("Done with  " + file.Name())
		}
	}
	wg.Wait()

	// Create blog root HTML
	newpath := filepath.Join(".", "public")
	os.MkdirAll(newpath, os.ModePerm)
	f, err := os.Create("public/index.html")
	if err != nil {
		log.Fatal(err)
	}
	t, err := template.ParseFiles("index.html")
	if err != nil {
		log.Fatalf("Failed to parse: %s", err)
	}
	if err := t.Execute(f, posts); err != nil {
		log.Fatalf("Failed to execute template: %s", err)
	}
	f.Close()

	// Copy styles
	if err := copyCSSToPublicDir(); err != nil {
		log.Fatalf("Failed to copy css: %s", err)
	}

	elapsed := time.Since(start)
	log.Printf("Done in %s", elapsed)
}

func createPost(filename string) (*Post, error) {
	// Get filename without extension
	filenameBase := fnameWithoutExtension(filename)
	if err := verifyFilenameBaseFormat(filenameBase); err != nil {
		return nil, err
	}

	// Get date and slug from filename
	date, err := time.Parse(dateLayout, filenameBase[0:len(dateLayout)])
	if err != nil {
		return nil, err
	}
	slug := filenameBase[11:]

	// Get body from file
	mdfile, err := os.Open("./markdown/" + filename)
	if err != nil {
		return nil, err
	}
	defer mdfile.Close()

	rawBytes, err := ioutil.ReadAll(mdfile)

	// Get title from first line of file
	lines := strings.Split(string(rawBytes), "\n")
	title := strings.Replace(lines[0], "# ", "", -1)

	// Convert Markdown to HTML
	html := blackfriday.Run(rawBytes)

	// Create a Post struct
	return &Post{
		Title: title,
		Body:  template.HTML(html),
		Date:  date,
		Slug:  slug,
	}, nil
}

func createPostFile(post *Post) error {
	// Create folder for HTML
	newpath := filepath.Join("public/posts", post.Slug)
	if err := os.MkdirAll(newpath, os.ModePerm); err != nil {
		return err
	}

	// Create HTML file
	f, err := os.Create("public/posts/" + post.Slug + "/" + "index.html")
	if err != nil {
		return err
	}
	defer f.Close()

	// Generate final HTML file from template
	t, err := template.ParseFiles("post.html")
	if err != nil {
		return err
	}
	return t.Execute(f, post)
}

func fnameWithoutExtension(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}

func verifyFilenameBaseFormat(f string) error {
	filenameRequirements := fmt.Sprintf(
		"Make sure file name is formatted as %s-slug.md",
		dateLayout)

	if len(f) < 12 {
		return fmt.Errorf("The name of the file is too short. %s", filenameRequirements)
	}

	_, err := time.Parse(dateLayout, f[0:len(dateLayout)])
	if err != nil {
		return fmt.Errorf("Date of the file doesn't look right. %s", filenameRequirements)
	}
	return nil
}

func copyCSSToPublicDir() error {
	from, err := os.Open("./css/styles.css")
	if err != nil {
		return err
	}
	defer from.Close()

	newpath := filepath.Join("public", "css")
	if err := os.MkdirAll(newpath, os.ModePerm); err != nil {
		return err
	}

	to, err := os.OpenFile("./public/css/styles.css", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	return err
}

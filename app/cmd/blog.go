package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/freetonik/underblog/app/internal"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"
)

const DefaultMarkdownPath = "./markdown/"

// Create and initialize Blog
func NewBlog(opts internal.Opts) *Blog {
	b := new(Blog)

	b.opts = opts

	b.mux = &sync.Mutex{}
	b.wg = &sync.WaitGroup{}

	b.files = make(chan os.FileInfo)

	return b
}

// Blog options and blog creating methods
type Blog struct {
	opts internal.Opts

	files     chan os.FileInfo
	posts     []Post
	indexPage io.Writer

	mux *sync.Mutex
	wg *sync.WaitGroup
}

// Render md-files->HTML, generate root index.html
func (b *Blog) Render() error {
	if err := b.verifyMarkdownPresent(); err != nil {
		log.Fatal(errors.New(fmt.Sprintf("Markdown directory is not found: %v", err)))
	}

	b.indexPage = b.getIndexPage(b.opts.Path)
	b.createPosts()
	err := b.renderMd()
	b.copyCssToPublicDir()

	return err
}

func (b *Blog) addPost(post Post) {
	b.mux.Lock()
	b.posts = append(b.posts, post)
	b.mux.Unlock()
}

func (b *Blog) getIndexPage(currentPath string) io.Writer {
	rootPath := "."

	if currentPath != "" {
		rootPath = currentPath
	}
	p := filepath.Join(rootPath, "public")
	err := os.MkdirAll(p, os.ModePerm)
	if err != nil {
		log.Fatal(errors.New(fmt.Sprintf("Can't create public dir: %v", err)))
	}

	f, err := os.Create("public/index.html")

	if err != nil {
		log.Fatal(errors.New(fmt.Sprintf("Can't create public/index.html: %v", err)))
	}

	return f
}

func (b *Blog) startWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case file, ok := <-b.files:
			if !ok || !isFileValid(file) {
				// todo: catch it?
				b.wg.Done()
				return
			}
			b.addPost(NewPost(file.Name()))
			b.wg.Done()
		}
	}

}

func (b *Blog) getMdFiles() []os.FileInfo {
	files, err := ioutil.ReadDir(DefaultMarkdownPath)
	if err != nil {
		fmt.Println("Can't get directory of markdown files")
		log.Fatal(err)
	}
	return files
}

func (b *Blog) createPosts() {
	ctx := context.Background()

	filesChan := make(chan os.FileInfo)
	files := b.getMdFiles()

	wLimit := internal.GetWorkersLimit(len(files))
	b.wg.Add(len(files))

	for i := 0; i < wLimit; i++ {
		go b.startWorker(ctx)
	}

	for _, file := range files {
		b.files <- file
	}

	close(filesChan)
}

func (b *Blog) copyCssToPublicDir() {
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

func (b *Blog) renderMd() error {
	t, _ := template.ParseFiles("index.html")
	b.wg.Wait() // wait until b.posts is populated
	err := t.Execute(b.indexPage, b.posts)
	if err != nil {
		log.Fatalf("can't execute template: %v", err)
	}
	// todo: should i close file interface?
	return nil
}

func (b *Blog) verifyMarkdownPresent() error {
	if _, err := os.Stat(DefaultMarkdownPath); os.IsNotExist(err) {
		return err
	}
	return nil
}

func isFileValid(file os.FileInfo) bool {
	return path.Ext(file.Name()) == ".md" || path.Ext(file.Name()) == ".markdown"
}

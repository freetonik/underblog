package main

import (
	"fmt"
	"github.com/freetonik/underblog/app/cmd"
	"github.com/freetonik/underblog/app/internal"
	"log"
	"os"
	"time"
)

var revision = "0.2.2"

func main() {
	fmt.Printf("Underblog %s\n", revision)

	opts := internal.GetCLIOptions()

	if opts.Version {
		fmt.Printf("Current ver.: %s\n", revision)
		os.Exit(0)
	}

	if opts.WatchMode {
		fmt.Println("Starting in watch mode...")
		go makeBlog(opts)
	} else {
		fmt.Println("Starting...")
		makeBlog(opts)
	}

	if opts.WatchMode {
		go internal.WatchForChangedFiles(func() { makeBlog(opts) })
		internal.RunDevelopmentWebserver()
	}
}

func makeBlog(opts internal.Opts) {
	start := time.Now()
	err := cmd.MakeBlog(opts)
	if err != nil {
		log.Fatalf("Can't make a blog: %v", err)
	}
	elapsed := time.Since(start)
	log.Printf("Done in %s", elapsed)
}

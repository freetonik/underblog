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

	start := time.Now()

	fmt.Printf("Starting...\n")
	err := cmd.MakeBlog(opts)
	if err != nil {
		log.Fatalf("Can't make a blog: %v", err)
	}

	elapsed := time.Since(start)
	log.Printf("Done in %s", elapsed)
}

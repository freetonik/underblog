package internal

import (
	"flag"
)

// Opts options of the command line
type Opts struct {
	Version   bool
	WatchMode bool
	Path      string
}

// GetCLIOptions parse CLI args and make Opts
func GetCLIOptions() Opts {
	// todo:
	//  - [] add path to "./markdown/" as cli args,
	//    https://github.com/freetonik/underblog/issues/10

	version := flag.Bool("version", false, "prints current version")
	watchMode := flag.Bool("watch", false, "launches in watch mode")

	flag.Parse()

	return Opts{
		Version:   *version,
		WatchMode: *watchMode,
		Path:      "",
	}
}

package internal

import (
	"flag"
)

type Opts struct {
	Version bool
	Path    string
}

// GetCLIOptions parse CLI args and make Opts
func GetCLIOptions() Opts {
	// todo:
	//  - [] add path to "./markdown/" as cli args,
	//    https://github.com/freetonik/underblog/issues/10

	version := flag.Bool("version", false, "prints current version")
	flag.Parse()
	return Opts{
		Version: *version,
		Path:    "",
	}
}

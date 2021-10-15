package utils

import (
	"fmt"

	"github.com/docopt/docopt-go"
)

// Arguments holds the parsed command line options of the application.
type Arguments struct {
	DebugMode  bool   `docopt:"--debug"`
	ConfigPath string `docopt:"--configPath"`
}

// ParseArguments provides the docopt command line interface and returns the parsed arguments
// when appropriate.
func ParseArguments() (*Arguments, error) {
	usage := `crate.run API.

    Usage:
      crate-api [--configPath=<filename>]
      crate-api (-h | --help)
      crate-api --version
	  crate-api --debug

    Options:
      -h --help            Show this screen.
      --version            Show version.
	  --debug              Run in debug mode.
      --configPath=<path>  Path of configuration file [default: config.toml].`

	opts, err := docopt.ParseArgs(usage, nil, "crate.run API 1.1.1")

	if err != nil {
		return nil, fmt.Errorf("could not parse args: %w", err)
	}

	var arguments Arguments
	if err := opts.Bind(&arguments); err != nil {
		return nil, fmt.Errorf("could not bind args: %w", err)
	}

	return &arguments, nil
}

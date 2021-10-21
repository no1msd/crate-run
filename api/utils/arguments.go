package utils

import (
	"flag"
)

// Arguments holds the parsed command line options of the application.
type Arguments struct {
	DebugMode  bool
	ConfigPath string
}

// ParseArguments returns the parsed command line arguments.
func ParseArguments() Arguments {
	debugFlag := flag.Bool("debug", false, "run in debug mode")
	configPath := flag.String("configPath", "config.toml", "path of configuration file")

	flag.Parse()

	return Arguments{
		DebugMode:  *debugFlag,
		ConfigPath: *configPath,
	}
}

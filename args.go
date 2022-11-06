package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	Version bool `short:"v" long:"version" description:"Show version"`

	Color string `short:"c" long:"color" description:"Matrix colors, can be up to 2 comma-separated colors for gradient" default:"000000,00FF00"`
}

type Config struct {
	showVersion bool
	colors      Colors
}

func ParseArgs() Config {
	_, err := flags.Parse(&opts)

	if err != nil {
		flagError := err.(*flags.Error)
		if flagError.Type == flags.ErrHelp {
			os.Exit(0)
		} else if flagError.Type == flags.ErrUnknownFlag {
			fmt.Println("Use --help to view all available options.")
			os.Exit(0)
		} else {
			fmt.Printf("Error parsing flags: %s\n", err)
			os.Exit(1)
		}
	}

	return Config{
		showVersion: opts.Version,
		colors:      parseColors(opts.Color),
	}
}

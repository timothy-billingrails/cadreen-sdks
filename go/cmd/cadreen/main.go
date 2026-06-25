package main

import (
	"os"

	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/commands"
)

var (
	version = "0.2.5"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	commands.Version = version
	commands.Commit = commit
	commands.Date = date

	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}

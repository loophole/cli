package main

import (
	"github.com/loophole/cli/cmd"
	"github.com/loophole/cli/internal/pkg/closehandler"
)

// Will be filled in during build
var version = "development"
var commit = "unknown"

func main() {
	c := closehandler.SetupCloseHandler("https://forms.gle/K9ga7FZB3deaffnV7") //needs to be here instead of root to avoid duplicates
	cmd.Execute(version, commit, c)
}

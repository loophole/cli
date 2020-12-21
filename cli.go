// +build !desktop

package main

import (
	"github.com/loophole/cli/cmd"
	"github.com/loophole/cli/config"
	"github.com/loophole/cli/internal/pkg/closehandler"
)

var (
	version = "development"
	commit  = "unknown"
	mode    = "cli"
)

func main() {
	config.Config.Version = version
	config.Config.CommitHash = commit
	config.Config.ClientMode = mode

	c := closehandler.SetupCloseHandler("https://forms.gle/K9ga7FZB3deaffnV7")
	cmd.Execute(c)
}

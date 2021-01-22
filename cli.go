// +build !desktop

package main

import (
	"github.com/loophole/cli/cmd"
	"github.com/loophole/cli/config"
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

	cmd.Execute()
}

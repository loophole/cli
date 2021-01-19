// +build desktop

package main

import (
	"github.com/loophole/cli/config"
	"github.com/loophole/cli/ui"
)

var (
	version = "development"
	commit  = "unknown"
	mode    = "desktop"
)

func main() {
	config.Config.Version = version
	config.Config.CommitHash = commit
	config.Config.ClientMode = mode

	ui.Display()
}

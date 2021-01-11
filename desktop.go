// +build desktop

package main

import (
	"github.com/loophole/cli/config"
	"github.com/loophole/cli/ui"
)

func main() {
	ui.Display(config.Config)
}

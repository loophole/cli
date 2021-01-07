// +build !desktop

package main

import (
	"github.com/loophole/cli/cmd"
	"github.com/loophole/cli/config"
)

func main() {
	cmd.Execute(config.Config)
}

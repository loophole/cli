package main

import (
	"github.com/loophole/cli/cmd"
)

// Will be filled in during build
var version = "development"
var commit = "unknown"

func main() {
	cmd.Execute(version, commit)
}

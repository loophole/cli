package inpututil

import (
	"os"

	"github.com/loophole/cli/internal/pkg/communication"
)

//IsUsingPipe returns whether a pipe is being used for inputs or not
func IsUsingPipe() bool {
	stdinInfo, err := os.Stdin.Stat()
	if err != nil {
		communication.Warn("Error getting terminal info")
		communication.Fatal(err.Error())
	}
	return (stdinInfo.Mode() & os.ModeCharDevice) == 0
}

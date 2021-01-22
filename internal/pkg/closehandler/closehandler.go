package closehandler

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/loophole/cli/internal/pkg/communication"
	"golang.org/x/crypto/ssh/terminal"
)

var successfulConnectionOccured bool = false
var terminalState *terminal.State = &terminal.State{}

// SetupCloseHandler ensures that CTRL+C inputs are properly processed, restoring the terminal state from not displaying entered characters where necessary
func SetupCloseHandler(feedbackFormURL string) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	terminalState, err := terminal.GetState(int(os.Stdin.Fd()))
	if err != nil {
		communication.Warn("Error saving terminal state")
		communication.Warn(err.Error())
	}

	go func() {
		<-c
		if terminalState != nil {
			terminal.Restore(int(os.Stdin.Fd()), terminalState)
		}
		communication.ApplicationStop()
		os.Exit(0)
	}()
}

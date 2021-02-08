package closehandler

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/loophole/cli/internal/pkg/communication"
	"github.com/loophole/cli/internal/pkg/inpututil"
	"golang.org/x/term"
)

var successfulConnectionOccured bool = false
var terminalState *term.State = &term.State{}

// SetupCloseHandler ensures that CTRL+C inputs are properly processed, restoring the terminal state from not displaying entered characters where necessary
func SetupCloseHandler(feedbackFormURL string) {
	var terminalState *term.State
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	if !inpututil.IsUsingPipe() { //don't try to get terminal state if using a pipe
		var err error
		terminalState, err = term.GetState(int(os.Stdin.Fd()))
		if err != nil {
			communication.Warn("Error saving terminal state")
			communication.Fatal(err.Error())
		}
	}
	go func() {
		<-c
		if terminalState != nil {
			term.Restore(int(os.Stdin.Fd()), terminalState)
		}
		communication.ApplicationStop()
		os.Exit(0)
	}()
}

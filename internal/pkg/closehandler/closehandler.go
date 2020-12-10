package closehandler

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/loophole/cli/internal/pkg/communication"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh/terminal"
)

var successfulConnectionOccured bool = false
var terminalState *terminal.State = &terminal.State{}

//FeedbackFormURL contains the link to the feedbackform
var FeedbackFormURL string = "https://forms.gle/K9ga7FZB3deaffnV7"

//SetupCloseHandler ensures that CTRL+C inputs are properly processed, restoring the terminal state from not displaying entered characters where necessary
func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	terminalState, err := terminal.GetState(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal().Err(err).Msg("Error saving terminal state")
	}

	go func() {
		<-c
		if terminalState != nil {
			terminal.Restore(int(os.Stdin.Fd()), terminalState)
		}
		communication.PrintGoodbyeMessage()
		if successfulConnectionOccured {
			communication.PrintFeedbackMessage(FeedbackFormURL)
		}
		os.Exit(0)
	}()
}

//SuccessfulConnectionOccured sets the corresponding boolean to true, enabling the display of the feedback form URL after closing the CLI
func SuccessfulConnectionOccured() {
	successfulConnectionOccured = true
}

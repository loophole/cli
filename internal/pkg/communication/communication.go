package communication

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/logrusorgru/aurora"
	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog/log"
)

var colorableOutput = colorable.NewColorableStdout()
var loader = spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithWriter(colorableOutput))

func PrintWelcomeMessage() {
	fmt.Fprint(colorableOutput, aurora.Cyan("Loophole"))
	fmt.Fprint(colorableOutput, aurora.Italic(" - End to end TLS encrypted TCP communication between you and your clients"))
	fmt.Println()
	fmt.Println()
}

func PrintTunnelSuccessMessage(siteAddr string, localAddr string) {
	NewLine()
	fmt.Fprint(colorableOutput, "Forwarding ")
	fmt.Fprint(colorableOutput, aurora.Green(siteAddr))
	fmt.Fprint(colorableOutput, " -> ")
	fmt.Fprint(colorableOutput, aurora.Green(localAddr))
	NewLine()
	WriteItalic("TLS Certificate will be obtained with first request to the above address, therefore first execution may be slower")
	NewLine()
	WriteCyan("Press CTRL + C to stop the service")
	NewLine()
	Write("Logs:")

	log.Info().Msg("Awaiting connections...")
}

func PrintGoodbyeMessage() {
	NewLine()
	Write("Goodbye")
}

func PrintFeedbackMessage(feedbackFormURL string) {
	fmt.Println(aurora.Cyan(fmt.Sprintf("Thank you for using Loophole. Please give us your feedback via %s and help us improve our services.", feedbackFormURL)))
}

func StartLoading(message string) {
	if el := log.Debug(); !el.Enabled() {
		loader.Prefix = fmt.Sprintf("%s ", message)
		loader.Start()
	} else {
		Write(message)
	}
}

func LoadingSuccess() {
	if el := log.Debug(); !el.Enabled() {
		loader.FinalMSG = fmt.Sprintf("%s%s\n", loader.Prefix, aurora.Green("Success!"))
		loader.Stop()
	}
}

func LoadingFailure() {
	if el := log.Debug(); !el.Enabled() {
		loader.FinalMSG = fmt.Sprintf("%s%s\n", loader.Prefix, aurora.Red("Error!"))
		loader.Stop()
	}
}

func NewLine() {
	fmt.Fprintln(colorableOutput)
}

func Write(message string) {
	fmt.Fprint(colorableOutput, fmt.Sprintf("%s", message))
	NewLine()
}

func WriteRed(message string) {
	fmt.Fprint(colorableOutput, fmt.Sprintf("%s", aurora.Red(message)))
	NewLine()
}
func WriteGreen(message string) {
	fmt.Fprint(colorableOutput, fmt.Sprintf("%s", aurora.Green(message)))
	NewLine()
}

func WriteCyan(message string) {
	fmt.Fprint(colorableOutput, fmt.Sprintf("%s", aurora.Cyan(message)))
	NewLine()
}

func WriteItalic(message string) {
	fmt.Fprint(colorableOutput, fmt.Sprintf("%s", aurora.Italic(message)))
	NewLine()
}

package communication

import (
	"fmt"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/logrusorgru/aurora"
	"github.com/loophole/cli/internal/pkg/urlmaker"
	"github.com/mattn/go-colorable"
	"github.com/mdp/qrterminal"
	"github.com/rs/zerolog/log"
)

var colorableOutput = colorable.NewColorableStdout()
var loader = spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithWriter(colorableOutput))
var MessageMutex sync.Mutex

func PrintWelcomeMessage() {
	MessageMutex.Lock()
	fmt.Fprint(colorableOutput, aurora.Cyan("Loophole"))
	fmt.Fprint(colorableOutput, aurora.Italic(" - End to end TLS encrypted TCP communication between you and your clients"))
	NewLine()
	NewLine()
	MessageMutex.Unlock()
}

func PrintTunnelSuccessMessage(siteID string, protocols []string, localAddr string, displayQR bool) {
	MessageMutex.Lock()
	NewLine()

	if len(protocols) < 1 {
		protocols = []string{"https"}
	}

	for _, protocol := range protocols {
		NewLine()
		siteAddr := urlmaker.GetSiteUrl(protocol, siteID)
		fmt.Fprint(colorableOutput, "Forwarding ")
		fmt.Fprint(colorableOutput, aurora.Green(siteAddr))
		fmt.Fprint(colorableOutput, " -> ")
		fmt.Fprint(colorableOutput, aurora.Green(localAddr))
	}

	if displayQR {
		NewLine()
		NewLine()
		Write("Scan the below QR code to open the site:")
		NewLine()
		QRCode(urlmaker.GetSiteUrl(protocols[0], siteID))
	}

	if len(protocols) > 1 {
		NewLine()
		NewLine()
		fmt.Fprint(colorableOutput, "Choose the protocol suitable for your target OS and share it")
		NewLine()
	}

	NewLine()
	WriteCyan("Press CTRL + C to stop the service")
	NewLine()
	Write("Logs:")

	log.Info().Msg("Awaiting connections...")
	MessageMutex.Unlock()
}

func PrintGoodbyeMessage() {
	MessageMutex.Lock()
	NewLine()
	Write("Goodbye")
	MessageMutex.Unlock()
}

func PrintFeedbackMessage(feedbackFormURL string) {
	MessageMutex.Lock()
	fmt.Fprintln(colorableOutput, aurora.Cyan(fmt.Sprintf("Thank you for using Loophole. Please give us your feedback via %s and help us improve our services.", feedbackFormURL)))
	MessageMutex.Unlock()
}

func StartLoading(message string) {
	if el := log.Debug(); !el.Enabled() {
		loader.Prefix = fmt.Sprintf("%s ", message)
		loader.Start()
	} else {
		MessageMutex.Lock()
		Write(message)
		MessageMutex.Unlock()
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

func LogInfo(message string) {
	MessageMutex.Lock()
	log.Info().Msg(message)
	MessageMutex.Unlock()
}

func LogFatalErr(message string, err error) {
	MessageMutex.Lock()
	log.Fatal().Err(err).Msg(message)
	MessageMutex.Unlock()
}

func LogFatalMsg(message string) {
	MessageMutex.Lock()
	log.Fatal().Msg(message)
	MessageMutex.Unlock()
}

func LogDebug(message string) {
	MessageMutex.Lock()
	log.Debug().Msg(message)
	MessageMutex.Unlock()
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

func QRCode(siteAddr string) {
	QRconfig := qrterminal.Config{
		Level:     qrterminal.L,
		Writer:    colorableOutput,
		BlackChar: qrterminal.WHITE,
		WhiteChar: qrterminal.BLACK,
		QuietZone: 1,
	}
	qrterminal.GenerateWithConfig(siteAddr, QRconfig)
}

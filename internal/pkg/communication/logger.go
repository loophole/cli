package communication

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/logrusorgru/aurora"
	tm "github.com/loophole/cli/internal/pkg/token/models"
	"github.com/loophole/cli/internal/pkg/urlmaker"
	"github.com/mattn/go-colorable"
	"github.com/mdp/qrterminal/v3"
	"github.com/rs/zerolog/log"
)

type stdoutLogger struct {
	colorableOutput io.Writer
	loader          *spinner.Spinner
	messageMutex    sync.Mutex
}

func NewStdOutLogger() Mechanism {
	logger := stdoutLogger{
		colorableOutput: colorable.NewColorableStdout(),
	}

	logger.loader = spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithWriter(logger.colorableOutput))

	return &logger
}

func (l *stdoutLogger) PrintLoginPrompt(deviceCodeSpec tm.DeviceCodeSpec) {
	fmt.Fprintf(l.colorableOutput, "Please open %s and use %s code to log in", aurora.Yellow(deviceCodeSpec.VerificationURI), aurora.Yellow(deviceCodeSpec.UserCode))
}

func (l *stdoutLogger) PrintWelcomeMessage(loggedIn bool) {
	l.messageMutex.Lock()
	fmt.Fprint(l.colorableOutput, aurora.Cyan("Loophole"))
	fmt.Fprint(l.colorableOutput, aurora.Italic(" - End to end TLS encrypted TCP communication between you and your clients"))
	l.NewLine()
	l.NewLine()
	l.messageMutex.Unlock()
}

func (l *stdoutLogger) PrintTunnelSuccessMessage(siteID string, protocols []string, localAddr string, displayQR bool) {
	l.messageMutex.Lock()
	l.NewLine()

	if len(protocols) < 1 {
		protocols = []string{"https"}
	}

	for _, protocol := range protocols {
		l.NewLine()
		siteAddr := urlmaker.GetSiteURL(protocol, siteID)
		fmt.Fprint(l.colorableOutput, "Forwarding ")
		fmt.Fprint(l.colorableOutput, aurora.Green(siteAddr))
		fmt.Fprint(l.colorableOutput, " -> ")
		fmt.Fprint(l.colorableOutput, aurora.Green(localAddr))
	}

	if displayQR {
		l.NewLine()
		l.NewLine()
		l.Write("Scan the below QR code to open the site:")
		l.NewLine()
		l.QRCode(urlmaker.GetSiteURL(protocols[0], siteID))
	}

	if len(protocols) > 1 {
		l.NewLine()
		l.NewLine()
		fmt.Fprint(l.colorableOutput, "Choose the protocol suitable for your target OS and share it")
		l.NewLine()
	}

	l.NewLine()
	l.WriteCyan("Press CTRL + C to stop the service")
	l.NewLine()
	l.Write("Logs:")

	log.Info().Msg("Awaiting connections...")
	l.messageMutex.Unlock()
}

func (l *stdoutLogger) PrintGoodbyeMessage() {
	l.messageMutex.Lock()
	l.NewLine()
	l.Write("Goodbye")
	l.messageMutex.Unlock()
}

func (l *stdoutLogger) PrintFeedbackMessage(feedbackFormURL string) {
	l.messageMutex.Lock()
	fmt.Fprintln(l.colorableOutput, aurora.Cyan(fmt.Sprintf("Thank you for using Loophole. Please give us your feedback via %s and help us improve our services.", feedbackFormURL)))
	l.messageMutex.Unlock()
}

func (l *stdoutLogger) StartLoading(message string) {
	if el := log.Debug(); !el.Enabled() {
		l.loader.Prefix = fmt.Sprintf("%s ", message)
		l.loader.Start()
	} else {
		l.messageMutex.Lock()
		l.Write(message)
		l.messageMutex.Unlock()
	}
}

func (l *stdoutLogger) LoadingSuccess() {
	if el := log.Debug(); !el.Enabled() {
		l.loader.FinalMSG = fmt.Sprintf("%s%s\n", l.loader.Prefix, aurora.Green("Success!"))
		l.loader.Stop()
	}
}

func (l *stdoutLogger) LoadingFailure() {
	if el := log.Debug(); !el.Enabled() {
		l.loader.FinalMSG = fmt.Sprintf("%s%s\n", l.loader.Prefix, aurora.Red("Error!"))
		l.loader.Stop()
	}
}

func (l *stdoutLogger) LogInfo(message string) {
	l.messageMutex.Lock()
	log.Info().Msg(message)
	l.messageMutex.Unlock()
}

func (l *stdoutLogger) LogWarnErr(message string, err error) {
	l.messageMutex.Lock()
	log.Warn().Err(err).Msg(message)
	l.messageMutex.Unlock()
}

func (l *stdoutLogger) LogFatalErr(message string, err error) {
	l.messageMutex.Lock()
	log.Fatal().Err(err).Msg(message)
	l.messageMutex.Unlock()
}

func (l *stdoutLogger) LogFatalMsg(message string) {
	l.messageMutex.Lock()
	log.Fatal().Msg(message)
	l.messageMutex.Unlock()
}

func (l *stdoutLogger) LogDebug(message string) {
	l.messageMutex.Lock()
	log.Debug().Msg(message)
	l.messageMutex.Unlock()
}

func (l *stdoutLogger) NewLine() {
	fmt.Fprintln(l.colorableOutput)
}

func (l *stdoutLogger) Write(message string) {
	fmt.Fprint(l.colorableOutput, fmt.Sprintf("%s", message))
	l.NewLine()
}

func (l *stdoutLogger) WriteRed(message string) {
	fmt.Fprint(l.colorableOutput, fmt.Sprintf("%s", aurora.Red(message)))
	l.NewLine()
}
func (l *stdoutLogger) WriteGreen(message string) {
	fmt.Fprint(l.colorableOutput, fmt.Sprintf("%s", aurora.Green(message)))
	l.NewLine()
}

func (l *stdoutLogger) WriteCyan(message string) {
	fmt.Fprint(l.colorableOutput, fmt.Sprintf("%s", aurora.Cyan(message)))
	l.NewLine()
}

func (l *stdoutLogger) WriteItalic(message string) {
	fmt.Fprint(l.colorableOutput, fmt.Sprintf("%s", aurora.Italic(message)))
	l.NewLine()
}

func (l *stdoutLogger) QRCode(siteAddr string) {
	qrterminal.GenerateHalfBlock(siteAddr, qrterminal.L, l.colorableOutput)
}

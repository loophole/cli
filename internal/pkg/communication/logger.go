package communication

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/logrusorgru/aurora"
	"github.com/loophole/cli/config"
	coreModels "github.com/loophole/cli/internal/app/loophole/models"
	authModels "github.com/loophole/cli/internal/pkg/token/models"
	"github.com/loophole/cli/internal/pkg/urlmaker"
	"github.com/mattn/go-colorable"
	"github.com/mdp/qrterminal"
	"github.com/rs/zerolog/log"
)

type stdoutLogger struct {
	colorableOutput io.Writer
	loader          *spinner.Spinner
	messageMutex    sync.Mutex
}

// NewStdOutLogger is stdout mechanism constructor
func NewStdOutLogger() Mechanism {
	logger := stdoutLogger{
		colorableOutput: colorable.NewColorableStdout(),
	}

	logger.loader = spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithWriter(logger.colorableOutput))

	return &logger
}

func (l *stdoutLogger) TunnelDebug(tunnelID string, message string) {
	l.messageMutex.Lock()
	defer l.messageMutex.Unlock()
	if el := log.Debug(); el.Enabled() {
		fmt.Println()
		el.Msg(message)
	}
}
func (l *stdoutLogger) TunnelInfo(tunnelID string, message string) {
	l.messageMutex.Lock()
	defer l.messageMutex.Unlock()
	log.Info().Msg(message)
}
func (l *stdoutLogger) TunnelWarn(tunnelID string, message string) {
	l.messageMutex.Lock()
	defer l.messageMutex.Unlock()
	log.Warn().Msg(message)
}
func (l *stdoutLogger) TunnelError(tunnelID string, message string) {
	l.messageMutex.Lock()
	defer l.messageMutex.Unlock()
	log.Error().Msg(message)
}

func (l *stdoutLogger) Debug(message string) {
	l.messageMutex.Lock()
	defer l.messageMutex.Unlock()
	if el := log.Debug(); el.Enabled() {
		fmt.Println()
		el.Msg(message)
	}
}
func (l *stdoutLogger) Info(message string) {
	l.messageMutex.Lock()
	defer l.messageMutex.Unlock()
	log.Info().Msg(message)
}
func (l *stdoutLogger) Warn(message string) {
	l.messageMutex.Lock()
	defer l.messageMutex.Unlock()
	log.Warn().Msg(message)
}
func (l *stdoutLogger) Error(message string) {
	l.messageMutex.Lock()
	defer l.messageMutex.Unlock()
	log.Error().Msg(message)
}
func (l *stdoutLogger) Fatal(message string) {
	l.messageMutex.Lock()
	defer l.messageMutex.Unlock()
	log.Fatal().Msg(message)
}

func (l *stdoutLogger) ApplicationStart(loggedIn bool, idToken string) {
	l.messageMutex.Lock()
	defer l.messageMutex.Unlock()
	fmt.Fprint(l.colorableOutput, aurora.Cyan("Loophole"))
	fmt.Fprint(l.colorableOutput, aurora.Italic(" - End to end TLS encrypted TCP communication between you and your clients"))
	fmt.Fprintln(l.colorableOutput)
	fmt.Fprintln(l.colorableOutput)
}
func (l *stdoutLogger) ApplicationStop() {
	l.messageMutex.Lock()
	defer l.messageMutex.Unlock()
	l.divider()
	fmt.Fprint(l.colorableOutput, "Goodbye")

	fmt.Fprintln(l.colorableOutput, aurora.Cyan(fmt.Sprintf("Thank you for using Loophole. Please give us your feedback via %s and help us improve our services.", config.Config.FeedbackFormURL)))
}

func (l *stdoutLogger) TunnelStart(tunnelID string) {
	l.messageMutex.Lock()
	defer l.messageMutex.Unlock()
	log.Debug().Msg("Tunnel starting up...")
}

func (l *stdoutLogger) TunnelStartSuccess(remoteConfig coreModels.RemoteEndpointSpecs, localEndpoint string) {
	l.messageMutex.Lock()
	defer l.messageMutex.Unlock()

	fmt.Fprintln(l.colorableOutput)
	fmt.Fprintln(l.colorableOutput)
	siteAddr := urlmaker.GetSiteURL("https", remoteConfig.SiteID, remoteConfig.Domain)
	fmt.Fprint(l.colorableOutput, "Forwarding ")
	fmt.Fprint(l.colorableOutput, aurora.Green(siteAddr))
	fmt.Fprint(l.colorableOutput, " -> ")
	fmt.Fprint(l.colorableOutput, aurora.Green(localEndpoint))

	if config.Config.Display.QR {
		fmt.Fprintln(l.colorableOutput, "")
		fmt.Fprintln(l.colorableOutput, "")

		log.Info().Msg("Scan the below QR code to open the site:")
		fmt.Fprintln(l.colorableOutput, "")
		qrterminal.GenerateHalfBlock(siteAddr, qrterminal.L, l.colorableOutput)
	}

	fmt.Fprintln(l.colorableOutput)
	fmt.Fprint(l.colorableOutput, aurora.Cyan("Press CTRL + C to stop the service"))
	fmt.Fprintln(l.colorableOutput)
	fmt.Fprintln(l.colorableOutput, "Logs: ")

	log.Info().Msg("Awaiting connections...")
}
func (l *stdoutLogger) TunnelStartFailure(tunnelID string, err error) {
	l.messageMutex.Lock()
	defer l.messageMutex.Unlock()
	log.Fatal().Str("tunnelId", tunnelID).Err(err).Msg("Tunnel startup error")
}
func (l *stdoutLogger) TunnelStopSuccess(tunnelID string) {
	l.messageMutex.Lock()
	defer l.messageMutex.Unlock()
	log.Debug().Str("tunnelId", tunnelID).Msg("Tunnel shutdown")
}

func (l *stdoutLogger) LoginStart(deviceCodeSpec authModels.DeviceCodeSpec) {
	l.messageMutex.Lock()
	defer l.messageMutex.Unlock()
	fmt.Fprintf(l.colorableOutput, "Please open %s and use %s code to log in", aurora.Yellow(deviceCodeSpec.VerificationURI), aurora.Yellow(deviceCodeSpec.UserCode))
}
func (l *stdoutLogger) LoginSuccess(idToken string) {
	l.messageMutex.Lock()
	defer l.messageMutex.Unlock()
	log.Info().Msg("Logged in successfully")
}
func (l *stdoutLogger) LoginFailure(err error) {
	l.messageMutex.Lock()
	defer l.messageMutex.Unlock()
	log.Fatal().Msg(err.Error())
}
func (l *stdoutLogger) LogoutSuccess() {
	l.messageMutex.Lock()
	defer l.messageMutex.Unlock()
	log.Info().Msg("Logged out successfully")
}
func (l *stdoutLogger) LogoutFailure(err error) {
	l.messageMutex.Lock()
	defer l.messageMutex.Unlock()
	log.Fatal().Msg(err.Error())
}

func (l *stdoutLogger) LoadingStart(tunnelID string, loaderMessage string) {
	l.messageMutex.Lock()
	defer l.messageMutex.Unlock()
	if el := log.Debug(); !el.Enabled() {
		l.loader.Prefix = fmt.Sprintf("%s ", loaderMessage)
		l.loader.Start()
	} else {
		fmt.Fprint(l.colorableOutput, loaderMessage)
	}
}

func (l *stdoutLogger) LoadingSuccess(tunnelID string) {
	l.messageMutex.Lock()
	defer l.messageMutex.Unlock()
	if el := log.Debug(); !el.Enabled() {
		l.loader.FinalMSG = fmt.Sprintf("%s%s\n", l.loader.Prefix, aurora.Green("Success!"))
		l.loader.Stop()
	}
}

func (l *stdoutLogger) LoadingFailure(tunnelID string, err error) {
	l.messageMutex.Lock()
	defer l.messageMutex.Unlock()
	if el := log.Debug(); !el.Enabled() {
		l.loader.FinalMSG = fmt.Sprintf("%s%s %s\n", l.loader.Prefix, aurora.Red("Error!"), err.Error())
		l.loader.Stop()
	}
}

func (l *stdoutLogger) NewVersionAvailable(availableVersion string) {
	l.messageMutex.Lock()
	defer l.messageMutex.Unlock()
	fmt.Fprint(l.colorableOutput, aurora.Cyan(fmt.Sprintf("There is new version available, to get it please visit %s",
		fmt.Sprintf("https://github.com/loophole/cli/releases/tag/%s", availableVersion))))
}

func (l *stdoutLogger) divider() {
	// No lock on purpose, becuase it's used internally
	fmt.Fprintln(l.colorableOutput)
}

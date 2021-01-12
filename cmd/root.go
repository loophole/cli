package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	stdlog "log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/AlecAivazis/survey/v2"
	lm "github.com/loophole/cli/internal/app/loophole/models"
	"github.com/loophole/cli/internal/pkg/cache"
	"github.com/loophole/cli/internal/pkg/closehandler"
	"github.com/loophole/cli/internal/pkg/communication"
	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

var displayOptions lm.DisplayOptions

var signalChan chan os.Signal

var alreadyRunning bool

//Possible answers for prompts and error messages
const (
	AnswerTunnelTypeHTTP   string = "Expose an HTTP Port"
	AnswerTunnelTypePath   string = "Expose a local path"
	AnswerTunnelTypeWebDAV string = "Expose a local path with WebDAV"
	AnswerYes              string = "Yes"
	AnswerNo               string = "No"
	PortRangeErrorMsg      string = "port must be between 0-65535"
	PathValidityErrorMsg   string = "enter an existing path without any quotation marks"
)

var rootCmd = &cobra.Command{
	Use:   "loophole",
	Short: "Loophole - End to end TLS encrypted TCP communication between you and your clients",
	Long:  "Loophole - End to end TLS encrypted TCP communication between you and your clients",
	Run: func(cmd *cobra.Command, args []string) {
		if !alreadyRunning {
			alreadyRunning = true
			interactivePrompt()
		}
	},
}

func getPortPrompt() []*survey.Question {
	return []*survey.Question{
		{
			Name:   "port",
			Prompt: &survey.Input{Message: "Please enter the http port you want to expose: "},
			Validate: func(val interface{}) error {
				if port, ok := val.(string); !ok {
					return errors.New(PortRangeErrorMsg)
				} else { //else is necessary here to keep access to port
					n, err := strconv.Atoi(port)
					if err != nil {
						return errors.New(PortRangeErrorMsg)
					}
					if (n < 0) || (n > 65535) {
						return errors.New(PortRangeErrorMsg)
					}
				}

				return nil
			},
		},
	}
}

func getPathPrompt() []*survey.Question {
	return []*survey.Question{
		{
			Name:   "path",
			Prompt: &survey.Input{Message: "Please enter the path you want to expose: "},
			Validate: func(val interface{}) error {
				if path, ok := val.(string); !ok {
					return errors.New(PathValidityErrorMsg)
				} else { //else is necessary here to keep access to path
					_, err := os.Stat(path)
					if err == nil {
						return nil
					}
					return errors.New(PathValidityErrorMsg)
				}
			},
		},
	}
}

func getLastArgsPrompt(lastArgs string) *survey.Select {
	return &survey.Select{
		Message: fmt.Sprintf("Your last settings were: '%s', would you like to reuse them?", lastArgs),
		Options: []string{AnswerYes, AnswerNo},
	}
}

func getInitialPrompt() *survey.Select {
	return &survey.Select{
		Message: "Welcome to loophole. What do you want to do?",
		Options: []string{AnswerTunnelTypeHTTP, AnswerTunnelTypePath, AnswerTunnelTypeWebDAV},
	}
}

func askBasicAuth() string {
	res := ""
	prompt := &survey.Select{
		Message: "Do you want to secure your tunnel using a username and password?",
		Options: []string{AnswerNo, AnswerYes},
	}
	var usernamePrompt = []*survey.Question{
		{
			Name:   "username",
			Prompt: &survey.Input{Message: "Please enter the username you want to use: "}, //not asking for a password since it's already implemented in virtual-serve
		},
	}
	err := survey.AskOne(prompt, &res)
	if err != nil {
		signalChan <- nil
	}
	if res == AnswerYes {
		err = survey.Ask(usernamePrompt, &res)
		if err != nil {
			os.Exit(0)
			return err.Error()
		}
	} else {
		return ""
	}
	return res
}

func askHostname() string {
	res := ""
	prompt := &survey.Select{
		Message: "Do you want to use a custom hostname?",
		Options: []string{AnswerNo, AnswerYes},
	}
	var hostnamePrompt = []*survey.Question{
		{
			Name:   "hostname",
			Prompt: &survey.Input{Message: "Please enter the hostname you want to use: "},
			Validate: func(val interface{}) error {
				var validChars = regexp.MustCompile(`^[a-z][a-z0-9]{0,30}$`).MatchString
				if hostname, ok := val.(string); !ok || len(hostname) > 31 || !validChars(hostname) || !unicode.IsLetter(rune(hostname[0])) {
					return errors.New("hostname must be up to 31 characters, may only contain lowercase letters and numbers and must start with a letter")
				}

				return nil
			},
		},
	}
	err := survey.AskOne(prompt, &res)
	if err != nil {
		signalChan <- nil
	}
	if res == AnswerYes {
		err = survey.Ask(hostnamePrompt, &res)
		if err != nil {
			os.Exit(0)
			return err.Error()
		}
	} else {
		return ""
	}
	return res
}

func interactivePrompt() {
	argPath := cache.GetLocalStorageFile("lastArgs", "logs")
	var lastArgs string = ""
	if _, err := os.Stat(argPath); err == nil {
		argBytes, err := ioutil.ReadFile(argPath)
		if err != nil {
			communication.LogFatalErr("Error reading last used arguments:", err)
		}
		lastArgs = string(argBytes)
	}
	var lastArgsPrompt = getLastArgsPrompt(lastArgs)
	var initialPrompt = getInitialPrompt()
	var portPrompt = getPortPrompt()
	var pathPrompt = getPathPrompt()

	var res string
	var exposePort int
	var exposePath string
	var arguments []string

	cmd := httpCmd.Root() //TODO: find a better way to access rootCMD

	if lastArgs != "" {
		err := survey.AskOne(lastArgsPrompt, &res)
		if err != nil {
			signalChan <- nil
		}
		if res == AnswerYes {
			cmd.SetArgs(strings.Split(lastArgs, " ")) //needs validation
			cmd.Execute()
			os.Exit(1)
		}
	}
	err := survey.AskOne(initialPrompt, &res)

	if err != nil {
		signalChan <- nil
	}
	switch res {
	case AnswerTunnelTypeHTTP:
		err = survey.Ask(portPrompt, &exposePort)
		if err != nil {
			signalChan <- nil
		}
		arguments = []string{"http", strconv.Itoa(exposePort)}
	case AnswerTunnelTypePath:
		err = survey.Ask(pathPrompt, &exposePath)
		if err != nil {
			signalChan <- nil
		}
		arguments = []string{"path", exposePath}
	case AnswerTunnelTypeWebDAV:
		err = survey.Ask(pathPrompt, &exposePath)
		if err != nil {
			signalChan <- nil
		}
		arguments = []string{"webdav", exposePath}
	}

	hostname := askHostname()
	if hostname != "" {
		arguments = append(arguments, "--hostname", hostname)
	}
	basicAuth := askBasicAuth()
	if basicAuth != "" {
		arguments = append(arguments, "--basic-auth-username", basicAuth)
	}
	cmd.SetArgs(arguments)

	var argumentsWithQuotes []string
	//setting the path argument in code doesn't work when it contains quotation marks,
	//but they do need to be there when entered as a standalone command in a command line if the path contains spaces
	//so, we give a copy of the arguments to the closehandler, with the path in quotation marks, where necessary
	if strings.Contains(exposePath, " ") {
		for i := 0; i < len(arguments); i++ {
			if arguments[i] == exposePath {
				argumentsWithQuotes = append(argumentsWithQuotes, fmt.Sprintf("'%s'", exposePath))
			} else {
				argumentsWithQuotes = append(argumentsWithQuotes, arguments[i])
			}
		}
		closehandler.SaveArguments(argumentsWithQuotes)
	} else {
		closehandler.SaveArguments(arguments)
	}
	cmd.Execute()
}

func init() {
	cobra.OnInitialize(initLogger)

	displayOptions.FeedbackFormURL = "https://forms.gle/K9ga7FZB3deaffnV7"
	rootCmd.PersistentFlags().BoolVarP(&displayOptions.Verbose, "verbose", "v", false, "verbose output")
}

func initLogger() {
	logLocation := cache.GetLocalStorageFile(fmt.Sprintf("%s.log", time.Now().Format("2006-01-02--15-04-05")), "logs")

	f, err := os.Create(logLocation)
	if err != nil {
		stdlog.Fatalln("Error creating log file:", err)
	}

	consoleWriter := zerolog.ConsoleWriter{Out: colorable.NewColorableStderr()}
	multi := zerolog.MultiLevelWriter(consoleWriter, f)
	log.Logger = zerolog.New(multi).With().Timestamp().Logger()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if displayOptions.Verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	stdlog.SetFlags(0)
	stdlog.SetOutput(log.Logger)
}

// Execute runs command parsing chain
func Execute(version string, commit string, c chan os.Signal) {
	rootCmd.Version = fmt.Sprintf("%s (%s)", version, commit)
	displayOptions.Version = fmt.Sprintf("%s-%s", version, commit)

	signalChan = c
	if !alreadyRunning {
		if err := rootCmd.Execute(); err != nil {
			signalChan <- nil
		}
	}
}

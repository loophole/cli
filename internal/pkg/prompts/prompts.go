package prompts

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/AlecAivazis/survey/v2"
	"github.com/loophole/cli/internal/pkg/cache"
	"github.com/loophole/cli/internal/pkg/closehandler"
	"github.com/loophole/cli/internal/pkg/communication"
	"github.com/spf13/cobra"
)

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

func askBasicAuth(signalChan chan os.Signal) string {
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

func askHostname(signalChan chan os.Signal) string {
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

func StartInteractivePrompt(cmd *cobra.Command, signalChan chan os.Signal) {
	argPath := cache.GetLocalStorageFile("lastArgs", "logs")
	var lastArgs string = ""
	if _, err := os.Stat(argPath); err == nil {
		argBytes, err := ioutil.ReadFile(argPath)
		if err != nil {
			communication.Fatal("Error reading last used arguments:" + err.Error())
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

	hostname := askHostname(signalChan)
	if hostname != "" {
		arguments = append(arguments, "--hostname", hostname)
	}
	basicAuth := askBasicAuth(signalChan)
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

// +build !desktop

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
	"github.com/loophole/cli/config"
	"github.com/loophole/cli/internal/pkg/cache"
	"github.com/loophole/cli/internal/pkg/closehandler"
	"github.com/loophole/cli/internal/pkg/token"
	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

var signalChan chan os.Signal

var alreadyRunning bool

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

func interactivePrompt() {
	cmd := httpCmd.Root() //find a better way to access rootCMD

	if !token.IsTokenSaved() {
		cmd.SetArgs([]string{"account", "login"})
		cmd.Execute()
	}

	argPath := cache.GetLocalStorageFile("lastArgs", "logs")
	var lastArgs string = ""
	if _, err := os.Stat(argPath); err == nil {
		argBytes, _ := ioutil.ReadFile(argPath)
		lastArgs = string(argBytes)
		fmt.Println(string(lastArgs))
	}
	argsq := &survey.Select{
		Message: "Your last settings were: " + lastArgs + " , would you like to reuse them?",
		Options: []string{"Yes", "No"},
	}
	initq := &survey.Select{
		Message: "Welcome to loophole. What do you want to do?",
		Options: []string{"Expose an HTTP Port", "Expose a local path", "Expose a local path with WebDAV", "Logout"},
	}
	var portPrompt = []*survey.Question{
		{
			Name:   "port",
			Prompt: &survey.Input{Message: "Please enter the http port you want to expose: "},
			Validate: func(val interface{}) error {
				if port, ok := val.(string); !ok {
					return errors.New("enter a valid string")
				} else { //else is necessary here to keep access to port
					n, err := strconv.Atoi(port)
					if err != nil {
						return errors.New("port must be between 0-65535")
					}
					if (n < 0) || (n > 65535) {
						return errors.New("port must be between 0-65535")
					}
				}

				return nil
			},
		},
	}
	var pathPrompt = []*survey.Question{
		{
			Name:   "path",
			Prompt: &survey.Input{Message: "Please enter the path you want to expose: "},
			Validate: func(val interface{}) error {
				if path, ok := val.(string); !ok {
					return errors.New("enter an existing path without any quotation marks")
				} else { //else is necessary here to keep access to path
					_, err := os.Stat(path)
					if err == nil {
						return nil
					}
					return errors.New("enter an existing path without any quotation marks")
				}
			},
			Transform: survey.TransformString(func(ans string) string {
				return fmt.Sprintf("'%s'", ans)
			}),
		},
	}
	logoutPrompt := &survey.Select{
		Message: "Are you sure you want to logout?",
		Options: []string{"No", "Yes, I'm sure"},
	}
	var res string
	var exposePort int
	var exposePath string
	var arguments []string

	if lastArgs != "" {
		err := survey.AskOne(argsq, &res)
		if err != nil {
			signalChan <- nil
		}
		if res == "Yes" {
			cmd.SetArgs(strings.Split(lastArgs, " ")) //needs validation
			cmd.Execute()
			os.Exit(0)
		}
	}
	err := survey.AskOne(initq, &res)
	if err != nil {
		signalChan <- nil
	}
	if res == "Expose an HTTP Port" {
		err = survey.Ask(portPrompt, &exposePort)
		if err != nil {
			signalChan <- nil
		}
		arguments = []string{"http", strconv.Itoa(exposePort)}
	} else if res == "Expose a local path" {
		err = survey.Ask(pathPrompt, &exposePath)
		if err != nil {
			signalChan <- nil
		}
		arguments = []string{"path", exposePath}
	} else if res == "Expose a local path with WebDAV" {
		err = survey.Ask(pathPrompt, &exposePath)
		if err != nil {
			signalChan <- nil
		}
		arguments = []string{"webdav", exposePath}
	} else if res == "Logout" {
		err := survey.AskOne(logoutPrompt, &res)
		if err != nil {
			signalChan <- nil
		}
		if res == "Yes, I'm sure" {
			cmd.SetArgs([]string{"account", "logout"})
			cmd.Execute()
		}
		os.Exit(0) //if Execute() should fail, don't ask for hostname etc. but instead exit
	}

	hostname := askHostname()
	if hostname != "" {
		arguments = append(arguments, "--hostname", hostname)
	}
	basicAuth := askBasicAuth()
	if basicAuth != "" {
		arguments = append(arguments, "-u", basicAuth)
	}
	closehandler.SaveArguments(arguments)
	cmd.SetArgs(arguments)
	cmd.Execute()
}

func askBasicAuth() string {
	res := ""
	prompt := &survey.Select{
		Message: "Do you want to secure your tunnel using a username and password?",
		Options: []string{"No", "Yes"},
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
	if res == "Yes" {
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
		Options: []string{"No", "Yes"},
	}
	var hostnamePrompt = []*survey.Question{
		{
			Name:   "hostname",
			Prompt: &survey.Input{Message: "Please enter the hostname you want to use: "},
			Validate: func(val interface{}) error {
				var validChars = regexp.MustCompile(`^[a-z0-9]+$`).MatchString
				if hostname, ok := val.(string); !ok || len(hostname) > 31 || len(hostname) < 6 || !validChars(hostname) || !unicode.IsLetter(rune(hostname[0])) {
					return errors.New("hostname must be between 6-31 characters, may only contain lowercase letters and numbers and must start with a letter")
				}

				return nil
			},
		},
	}
	err := survey.AskOne(prompt, &res)
	if err != nil {
		signalChan <- nil
	}
	if res == "Yes" {
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

func init() {
	cobra.OnInitialize(initLogger)

	rootCmd.PersistentFlags().BoolVarP(&config.Config.Display.Verbose, "verbose", "v", false, "verbose output")
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
	if config.Config.Display.Verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	stdlog.SetFlags(0)
	stdlog.SetOutput(log.Logger)
}

// Execute runs command parsing chain
func Execute(c chan os.Signal) {
	rootCmd.Version = fmt.Sprintf("%s (%s)", config.Config.Version, config.Config.CommitHash)

	signalChan = c
	if !alreadyRunning {
		if err := rootCmd.Execute(); err != nil {
			signalChan <- nil
		}
	}
}

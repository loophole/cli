// +build !desktop

package cmd

import (
	"fmt"
	stdlog "log"
	"os"
	"strconv"
	"strings"
	"time"

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

var b bool

var rootCmd = &cobra.Command{
	Use:   "loophole",
	Short: "Loophole - End to end TLS encrypted TCP communication between you and your clients",
	Long:  "Loophole - End to end TLS encrypted TCP communication between you and your clients",
	Run: func(cmd *cobra.Command, args []string) {
		if !b {
			b = true
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

	initq := &survey.Select{
		Message: "Welcome to loophole. What do you want to do?",
		Options: []string{"Expose an HTTP Port", "Expose a local path", "Expose a local path with WebDAV", "Logout"},
	}
	var portPrompt = []*survey.Question{
		{
			Name:      "port",
			Prompt:    &survey.Input{Message: "Please enter the http port you want to expose: "},
			Validate:  survey.Required,
			Transform: survey.Title,
		},
	}
	var pathPrompt = []*survey.Question{
		{
			Name:      "path",
			Prompt:    &survey.Input{Message: "Please enter the path you want to expose: "},
			Validate:  survey.Required,
			Transform: survey.Title,
		},
	}
	logoutPrompt := &survey.Select{
		Message: "Are you sure you want to logout?",
		Options: []string{"No", "Yes, I'm sure"},
	}
	var res string
	var exposePort int
	var exposePath string

	err := survey.AskOne(initq, &res)
	if err != nil {
		signalChan <- nil
	}
	if res == "Expose an HTTP Port" {
		err = survey.Ask(portPrompt, &exposePort)
		if err != nil {
			signalChan <- nil
		}
		hostname := askHostname()
		if hostname != "" {
			arguments := []string{"http", strconv.Itoa(exposePort), "--hostname", strings.ToLower(hostname)}
			closehandler.SaveArguments(arguments)
			cmd.SetArgs(arguments)
		} else {
			arguments := []string{"http", strconv.Itoa(exposePort)}
			closehandler.SaveArguments(arguments)
			cmd.SetArgs(arguments)
		}
		cmd.Execute()
	} else if res == "Expose a local path" {
		err = survey.Ask(pathPrompt, &exposePath)
		if err != nil {
			signalChan <- nil
		}
		hostname := askHostname()
		if hostname != "" {
			arguments := []string{"path", exposePath, "--hostname", strings.ToLower(hostname)}
			closehandler.SaveArguments(arguments)
			cmd.SetArgs(arguments)
		} else {
			arguments := []string{"path", exposePath}
			closehandler.SaveArguments(arguments)
			cmd.SetArgs(arguments)
		}
		cmd.Execute()
	} else if res == "Expose a local path with WebDAV" {
		err = survey.Ask(pathPrompt, &exposePath)
		if err != nil {
			signalChan <- nil
		}
		hostname := askHostname()
		if hostname != "" {
			arguments := []string{"webdav", exposePath, "--hostname", strings.ToLower(hostname)}
			closehandler.SaveArguments(arguments)
			cmd.SetArgs(arguments)
		} else {
			arguments := []string{"webdav", exposePath}
			closehandler.SaveArguments(arguments)
			cmd.SetArgs(arguments)
		}
		cmd.Execute()
	} else if res == "Logout" {
		err := survey.AskOne(logoutPrompt, &res)
		if err != nil {
			signalChan <- nil
		}
		if res == "Yes, I'm sure" {
			cmd.SetArgs([]string{"account", "logout"})
			cmd.Execute()
		}
	}
}

func askHostname() string {
	res := ""
	prompt := &survey.Select{
		Message: "Do you want to use a custom hostname?",
		Options: []string{"Yes", "No"},
	}
	var hostnamePrompt = []*survey.Question{
		{
			Name:      "hostname",
			Prompt:    &survey.Input{Message: "Please enter the hostname you want to use: "},
			Validate:  survey.Required,
			Transform: survey.Title,
		},
	}
	err := survey.AskOne(prompt, &res)
	if err != nil {
		signalChan <- nil
	}
	if res == "Yes" {
		err = survey.Ask(hostnamePrompt, &res)
		time.Sleep(1 * time.Second)
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
	if !b {
		if err := rootCmd.Execute(); err != nil {
			signalChan <- nil
		}
	}
}

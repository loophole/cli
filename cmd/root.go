package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	stdlog "log"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/loophole/cli/internal/app/loophole"
	lm "github.com/loophole/cli/internal/app/loophole/models"
	"github.com/mattn/go-colorable"
	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

var config lm.Config
var verbose bool

var rootCmd = &cobra.Command{
	Use:   "loophole <port> [host]",
	Short: "Loophole - End to end TLS encrypted TCP communication between you and your clients",
	Long:  "Loophole - End to end TLS encrypted TCP communication between you and your clients",
	Run: func(cmd *cobra.Command, args []string) {
		config.Host = "127.0.0.1"
		if len(args) > 1 {
			config.Host = args[1]
		}
		port, _ := strconv.ParseInt(args[0], 10, 32)
		config.Port = int32(port)
		loophole.Start(config)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Missing argument: port")
		}
		_, err := strconv.ParseInt(args[0], 10, 32)
		if err != nil {
			return fmt.Errorf("Invalid argument: port: %v", err)
		}
		return nil
	},
}

func init() {
	rootCmd.Version = "1.0.0"

	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	cobra.OnInitialize(initLogger)

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	rootCmd.Flags().StringVarP(&config.IdentityFile, "identity-file", "i", fmt.Sprintf("%s/.ssh/id_rsa", home), "private key path")
	rootCmd.Flags().StringVar(&config.GatewayEndpoint.Host, "gateway-url", "gateway.loophole.host", "remote gateway URL")
	rootCmd.Flags().Int32Var(&config.GatewayEndpoint.Port, "gateway-port", 8022, "remote gateway port")
	rootCmd.Flags().StringVar(&config.SiteID, "hostname", "", "custom hostname you want to run service on")
	rootCmd.Flags().BoolVar(&config.HTTPS, "https", false, "use if your local service is already using HTTPS")

}

func initLogger() {
	logLocation := "logs/" + time.Now().Format("2006-01-02--15-04-05") + ".log" //path to where the current log file will be saved at, named after the timestamp of its creation for now

	var err error

	if _, err := os.Stat("logs"); err != nil { //does the logs folder exist? If not, create it
		os.Mkdir("logs", 0700)
	}

	r, w, err := os.Pipe() //create a pipe that leads everything written into the writer back into the reader
	if err != nil {
		stdlog.Fatalf("Error creating pipe:" + err.Error())
	}

	go func(r *os.File) {
		var logMutex sync.Mutex //Mutex to prevent simultaneous read and write access to the log string

		buf := new(bytes.Buffer) //Buffer that will hold the unedited logs
		logstring := ""          //the log string that will be written into a file

		go func(buf *bytes.Buffer, r *os.File) { //Continuously read all logs into the Buffer
			for {
				time.Sleep(500 * time.Millisecond)
				_, err = io.Copy(buf, r)
				if err != nil {
					stdlog.Fatalf("Error creating log buffer: %v", err)
				}
			}
		}(buf, r)

		go func(logstring *string, logMutex *sync.Mutex) { //Continuously save an ANSI stripped version of those logs into a string
			for {
				time.Sleep(500 * time.Millisecond)
				exp := regexp.MustCompile("\\[\\d+m")
				exp2 := regexp.MustCompile("")
				logMutex.Lock()
				*logstring = exp2.ReplaceAllString((exp.ReplaceAllString(buf.String(), "")), "") //remove all occurrences of both regexes from the logs
				logMutex.Unlock()
			}
		}(&logstring, &logMutex)

		go func(logstring *string, logMutex *sync.Mutex) { //Continuously save that string into a file
			for {
				var logFile []byte //declaring the variable that will hold the data of the currently existing log file
				time.Sleep(500 * time.Millisecond)

				if _, err := os.Stat(logLocation); err == nil {
					logFile, _ = ioutil.ReadFile(logLocation)
				}
				logMutex.Lock()
				if !bytes.Equal([]byte(*logstring), logFile) { //check if there is something new to write to the logfile

					err := ioutil.WriteFile(logLocation, []byte(*logstring), 0777)
					if err != nil {
						fmt.Println(err)
					}
				}
				logMutex.Unlock()
			}
		}(&logstring, &logMutex)
	}(r)

	wrt := io.MultiWriter(colorable.NewColorableStderr(), w) // create a multiwriter that writes all it receives into both the ColorableStderr as well as the writer of the pipe
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: wrt}) // set the multiwriter to be the output of our logs
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	stdlog.SetFlags(0)
	stdlog.SetOutput(log.Logger)
}

// Execute runs command parsing chain
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

package loopholed

import (
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/takama/daemon"
)

var stdlog, errlog *log.Logger

// Service has embedded daemon
// daemon.Daemon abstracts os specific service mechanics
type Service struct {
	daemon.Daemon
}

func init() {
	stdlog = log.New(os.Stdout, "", 0)
	errlog = log.New(os.Stderr, "", 0)
}

// New return new daemon service
func New() *Service {
	daemonKind := daemon.SystemDaemon
	if runtime.GOOS == "darwin" {
		daemonKind = daemon.UserAgent
	}
	srv, err := daemon.New(name, description, daemonKind, dependencies...)
	if err != nil {
		errlog.Println("Error: ", err)
		os.Exit(1)
	}

	tmpl := srv.GetTemplate()
	err = srv.SetTemplate(strings.ReplaceAll(tmpl, "/var/run/", "/var/"))
	if err != nil {
		errlog.Println("Error: ", err)
		os.Exit(1)
	}
	return &Service{srv}
}

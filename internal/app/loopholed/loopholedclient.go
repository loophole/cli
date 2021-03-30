package loopholed

import (
	"fmt"
	"net"
	"os"

	"github.com/loophole/cli/internal/app/loophole/models"
	"github.com/loophole/cli/internal/pkg/communication"
)

type LoopholedClient struct {
}

func (client *LoopholedClient) Ps() {
	fmt.Println("Connecting...")
	conn, err := net.Dial("tcp", port)
	if err != nil {
		communication.Error("Cannot reach loophole daemon")
		os.Exit(1)
	}
	fmt.Println("Connected to the daemon")
	conn.Write([]byte("PS\n"))

	fmt.Println("Reading response...")
	for {
		buf := make([]byte, 4096)
		numbytes, err := conn.Read(buf)
		if numbytes == 0 || err != nil {
			return
		}
		fmt.Print(string(buf))
	}
}

func (client *LoopholedClient) HTTP(conf models.ExposeHTTPConfig) {
	fmt.Println("Connecting...")
	conn, err := net.Dial("tcp", port)
	if err != nil {
		communication.Error("Cannot reach loophole daemon")
		os.Exit(1)
	}
	fmt.Println("Connected to the daemon")
	message := fmt.Sprintf("HTTP,%d,%s,%s\n", conf.Local.Port, conf.Local.Host, conf.Remote.SiteID)
	conn.Write([]byte(message))

	fmt.Println("Reading response...")
	for {
		buf := make([]byte, 4096)
		numbytes, err := conn.Read(buf)
		if numbytes == 0 || err != nil {
			return
		}
		fmt.Print(string(buf))
	}
}

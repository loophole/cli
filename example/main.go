package main

import (
	"context"
	"fmt"
	"time"

	"github.com/loophole/cli/cmd"
)

func main() {
	fmt.Println("Master func")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("Started loophole tunnel to 127.0.0.1:80")
	go cmd.GoExecute(ctx, "development", "unknown", "cli", "--hostname", "slave", "http", "80")

	fmt.Println("Started local http server on 127.0.0.1:80")
	time.Sleep(time.Second * 30) // Mock server works 30s
	fmt.Println("Server stopped")
}

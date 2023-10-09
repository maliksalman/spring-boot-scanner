package main

import (
	"context"
	"github.com/maliksalman/spring-boot-scanner/cmd"
	"os"
	"os/signal"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	root := cmd.NewCmdRoot()
	if err := root.ExecuteContext(ctx); err != nil {
		cancel()
		os.Exit(1)
	}
}

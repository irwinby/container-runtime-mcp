package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/irwinby/container-runtime-mcp/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err := app.Run(ctx)
	if err != nil {
		_, _ = os.Stderr.WriteString("run app: " + err.Error() + "\n")
		os.Exit(1)
	}
}

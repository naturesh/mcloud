package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/naturesh/mcloud/internal/app"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	app.Execute(ctx)
}

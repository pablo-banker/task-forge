package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pablo-banker/taskforge/internal/app"
	"github.com/pablo-banker/taskforge/internal/config"
)

func main() {
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	application := app.New(cfg)

	if err := application.Run(ctx); err != nil {
		log.Fatal(err)
	}
}

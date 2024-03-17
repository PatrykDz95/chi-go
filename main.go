package main

import (
	"context"
	"fmt"
	"main/application"
	"os"
	"os/signal"
)

func main() {
	app := application.New()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	defer cancel()
	err := app.Start(ctx)
	if err != nil {
		fmt.Print("failed to start the server: %w", err)
	}
}

package main

import (
	"context"
	"flag"
	"log"

	"github.com/Slintox/user-service/internal/app"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "", "")
	flag.Parse()
}

func main() {
	ctx := context.Background()

	userApp, err := app.NewApp(ctx, configPath)
	if err != nil {
		log.Fatalf("failed to create the app: %s", err.Error())
	}

	if err = userApp.Run(); err != nil {
		log.Fatalf("failed to run the app: %s", err.Error())
	}
}

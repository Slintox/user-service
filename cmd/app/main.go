package main

import (
	"flag"

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
	app.Run(configPath)
}

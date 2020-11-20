package main

import (
	"os"

	"github.com/kushsharma/go-grpc-base/cmd"
	log "github.com/sirupsen/logrus"
)

var (
	// Version is app version
	Version = "0.0.1"
	// AppName of this executable
	AppName = "gbase"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	initConfig()

	rootCmd := cmd.InitCommands(AppName, Version)
	rootCmd.Execute()
}

func initConfig() {
	// TODO
}

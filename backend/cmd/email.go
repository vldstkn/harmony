package main

import (
	"harmony/internal/config"
	"harmony/internal/services/email"
	"harmony/pkg/logger"
	"os"
)

func main() {
	conf := config.LoadConfig()
	log := logger.NewLogger(os.Stdout)
	app := email.NewApp(&email.AppDeps{
		Config: conf,
		Logger: log,
	})
	app.Run()
}

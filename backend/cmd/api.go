package main

import (
	"harmony/internal/config"
	"harmony/internal/services/api"
	"harmony/pkg/logger"
	"os"
)

func main() {
	conf := config.LoadConfig()
	logger := logger.NewLogger(os.Stdout)
	app := api.NewApp(&api.AppDeps{
		Config: conf,
		Logger: logger,
	})
	app.Run()
}

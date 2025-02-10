package main

import (
	"harmony/internal/config"
	"harmony/internal/services/notifications"
	"harmony/pkg/logger"
	"os"
)

func main() {
	conf := config.LoadConfig()
	log := logger.NewLogger(os.Stdout)
	app := notifications.NewApp(&notifications.AppDeps{
		Logger: log,
		Config: conf,
	})
	app.Run()
}

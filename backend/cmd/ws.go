package main

import (
	"harmony/internal/config"
	"harmony/internal/services/ws"
	"harmony/pkg/logger"
	"os"
)

func main() {
	conf := config.LoadConfig()
	log := logger.NewLogger(os.Stdout)
	app := ws.NewApp(&ws.AppDeps{
		Logger: log,
		Config: conf,
	})
	app.Run()
}

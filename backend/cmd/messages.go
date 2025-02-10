package main

import (
	"harmony/internal/config"
	messages "harmony/internal/services/message"
	"harmony/pkg/db"
	"harmony/pkg/logger"
	"os"
)

func main() {
	conf := config.LoadConfig()
	database := db.NewDb(conf.DSN)
	log := logger.NewLogger(os.Stdout)
	app := messages.NewApp(&messages.AppDeps{
		Logger: log,
		Config: conf,
		DB:     database,
	})
	app.Run()
}

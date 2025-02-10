package main

import (
	"harmony/internal/config"
	"harmony/internal/services/room"
	"harmony/pkg/db"
	"harmony/pkg/logger"
	"os"
)

func main() {
	conf := config.LoadConfig()
	database := db.NewDb(conf.DSN)
	log := logger.NewLogger(os.Stdout)
	app := room.NewApp(&room.AppDeps{
		Config: conf,
		Logger: log,
		DB:     database,
	})
	app.Run()
}

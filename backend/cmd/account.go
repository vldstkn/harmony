package main

import (
	"harmony/internal/config"
	"harmony/internal/services/account"
	"harmony/pkg/db"
	"harmony/pkg/logger"
	"os"
)

func main() {
	conf := config.LoadConfig()
	log := logger.NewLogger(os.Stdout)
	database := db.NewDb(conf.DSN)
	app := account.NewApp(&account.AppDeps{
		Config: conf,
		Logger: log,
		DB:     database,
	})

	app.Run()
}

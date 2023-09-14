package main

import (
	"backend/cmd/job/wire"
	"backend/pkg/config"
	"backend/pkg/log"
)

func main() {
	conf := config.NewConfig()
	logger := log.NewLog(conf)
	logger.Info("start")

	app, cleanup, err := wire.NewApp(conf, logger)
	if err != nil {
		panic(err)
	}
	app.Run()
	defer cleanup()

}

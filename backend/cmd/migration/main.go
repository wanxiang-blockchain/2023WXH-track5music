package main

import (
	"backend/cmd/migration/wire"
	"backend/pkg/config"
	"backend/pkg/log"
)

func main() {
	conf := config.NewConfig()
	logger := log.NewLog(conf)

	app, cleanup, err := wire.NewApp(conf, logger)
	if err != nil {
		panic(err)
	}
	app.Run()
	defer cleanup()
}

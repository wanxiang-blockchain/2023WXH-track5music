package main

import (
	"backend/cmd/server/wire"
	"backend/pkg/config"
	"backend/pkg/http"
	"backend/pkg/log"
	"fmt"
	"go.uber.org/zap"
)

func main() {
	conf := config.NewConfig()
	logger := log.NewLog(conf)

	servers, cleanup, err := wire.NewApp(conf, logger)
	if err != nil {
		panic(err)
	}
	logger.Info("server start", zap.String("host", "http://localhost:"+conf.GetString("http.port")))

	//servers.
	http.Run(servers.ServerHTTP, fmt.Sprintf(":%d", conf.GetInt("http.port")))
	defer cleanup()

}

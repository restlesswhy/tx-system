package main

import (
	"fmt"
	"log"
	"txsystem/config"
	"txsystem/internal/server"
	"txsystem/pkg/logger"
	"txsystem/pkg/postgres"
)

func main() {
	log.Println("Starting microservice")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	pool, err := postgres.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()
	appLogger.Named(fmt.Sprintf(`(%s)`, cfg.ServiceName))
	appLogger.Infof("CFG: %+v", cfg)
	appLogger.Fatal(server.New(appLogger, cfg, pool).Run())
}

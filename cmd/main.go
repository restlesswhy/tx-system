package main

import (
	"txsystem/config"
	"txsystem/internal/server"
	"txsystem/pkg/postgres"

	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.Load()

	pool, err := postgres.Connect(cfg)
	if err != nil {
		logrus.Fatal(err)
	}

	srv := server.New(pool, cfg)
	srv.Run()
}

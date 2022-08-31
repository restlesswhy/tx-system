package main

import (
	"runtime"
	"txsystem/config"
	"txsystem/internal/server"
	"txsystem/pkg/postgres"

	"github.com/sirupsen/logrus"
)

func main() {
	runtime.SetBlockProfileRate(1)
	cfg := config.Load()

	pool, err := postgres.Connect(cfg)
	if err != nil {
		logrus.Fatal(err)
	}

	srv := server.New(pool, cfg)
	srv.Run()
}

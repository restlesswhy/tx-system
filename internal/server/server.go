package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"txsystem/config"
	"txsystem/internal/app"
	v1 "txsystem/internal/delivery/http/v1"
	"txsystem/internal/store"
	"txsystem/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"

	"net/http"
	_ "net/http/pprof"
)

type server struct {
	log   logger.Logger
	cfg   *config.Config
	pool  *pgxpool.Pool
	fiber *fiber.App
}

func New(log logger.Logger, cfg *config.Config, pool *pgxpool.Pool) *server {
	return &server{log: log, cfg: cfg, pool: pool, fiber: fiber.New()}
}

func (s *server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	store := store.New(s.pool)
	app := app.New(store)
	controller := v1.New(app)
	controller.SetupRoutes(s.fiber)

	go func() {
		if err := s.runHttp(); err != nil {
			s.log.Errorf("(runHttpServer) err: %v", err)
			cancel()
		}
	}()
	s.log.Infof("%s is listening on PORT: %v", s.getMicroserviceName(), s.cfg.Http.Port)

	go func() {
		s.log.Error(http.ListenAndServe(":6060", nil))
	}()

	<-ctx.Done()

	if err := s.fiber.Shutdown(); err != nil {
		s.log.Warnf("(Shutdown) err: %v", err)
	}
	app.Close()

	return nil
}

func (s server) getMicroserviceName() string {
	return fmt.Sprintf("(%s)", strings.ToUpper(s.cfg.ServiceName))
}

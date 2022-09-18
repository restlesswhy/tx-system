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
	"github.com/sirupsen/logrus"

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
	_, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	store := store.New(s.pool)
	app := app.New(store)
	controller := v1.New(app)
	s.fiber.Group(s.cfg.Http.BasePath)
	controller.SetupRoutes(s.fiber.Group(s.cfg.Http.BasePath))

	// r := mux.NewRouter()

	// srv := &http.Server{
	// 	Addr:    s.cfg.Addr,
	// 	Handler: v1.LoadRoutes(r, controller),
	// }

	go func() {
		if err := s.runHttp(); err != nil {
			s.log.Errorf("(runHttpServer) err: %v", err)
			cancel()
		}
	}()
	s.log.Infof("%s is listening on PORT: %v", s.getMicroserviceName(), s.cfg.Http.Port)

	// go func() {
	// 	if err := srv.ListenAndServe(); err != nil {
	// 		logrus.Fatal(err)
	// 	}
	// }()

	// go func() {
	// 	fmt.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch
	logrus.Info("stopping srv...")
	app.Close()

	return nil
}

func (s server) getMicroserviceName() string {
	return fmt.Sprintf("(%s)", strings.ToUpper(s.cfg.ServiceName))
}

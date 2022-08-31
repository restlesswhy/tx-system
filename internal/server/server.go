package server

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"txsystem/config"
	"txsystem/internal/app"
	"txsystem/internal/handler"
	"txsystem/internal/store"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"

	_ "net/http/pprof"
)

type Server struct {
	cfg  *config.Config
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool, cfg *config.Config) *Server {
	return &Server{pool: pool, cfg: cfg}
}

func (s *Server) Run() {
	store := store.New(s.pool)
	app := app.New(store)
	hndl := handler.New(app)

	r := mux.NewRouter()

	srv := &http.Server{
		Addr:    s.cfg.Addr,
		Handler: handler.LoadRoutes(r, hndl),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logrus.Fatal(err)
		}
	}()

	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	logrus.Infof("srv start listening on port %s...", s.cfg.Addr)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch
	logrus.Info("stopping srv...")
	app.Close()
}

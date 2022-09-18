package v1

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func (h *handler) SetupRoutes(r fiber.Router) {
	r.Use(logger.New())
	// api.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("PONG!")) })

	user := r.Group("/user")
	user.Post("/create", h.createUser)

	balance := r.Group("/balance")
	balance.Put("/update", h.updateBalance)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s\n", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

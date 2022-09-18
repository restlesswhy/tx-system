package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func (h *handler) SetupRoutes(r *fiber.App) {
	api := r.Group("/api/v1", logger.New())

	user := api.Group("/user")
	user.Post("/create", h.createUser)

	balance := api.Group("/balance")
	balance.Put("/update", h.updateBalance)
}

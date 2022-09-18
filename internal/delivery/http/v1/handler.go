package v1

import (
	"txsystem/internal/models"

	"github.com/gofiber/fiber/v2"
)

type App interface {
	ChangeBalance(tx *models.Transaction) error
	CreateUser(user *models.User) error
}

type handler struct {
	app App
}

func New(app App) *handler {
	return &handler{app: app}
}

func (h *handler) createUser(c *fiber.Ctx) error {
	user := &models.User{}
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
		})
	}

	if err := h.app.CreateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
	})
}

func (h *handler) updateBalance(c *fiber.Ctx) error {
	tx := &models.Transaction{}
	if err := c.BodyParser(tx); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
		})
	}

	if err := h.app.ChangeBalance(tx); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
	})
}

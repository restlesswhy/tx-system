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
			"result":  false,
			"message": err.Error(),
		})
	}

	if err := h.app.CreateUser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"result":  false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"result": true,
	})
}

func (h *handler) updateBalance(c *fiber.Ctx) error {
	tx := &models.Transaction{}

	if err := c.BodyParser(tx); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"result":  false,
			"message": err.Error(),
		})
	}

	if err := h.app.ChangeBalance(tx); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"result":  false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"result": true,
	})
}

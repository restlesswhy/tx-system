package server

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func (s *server) runHttp() error {
	s.fiber.Use(logger.New())
	s.fiber.Use(cors.New())

	return s.fiber.Listen(s.cfg.Http.Port)
}
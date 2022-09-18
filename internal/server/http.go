package server

func (s *server) runHttp() error {
	return s.fiber.Listen(s.cfg.Http.Port)
}
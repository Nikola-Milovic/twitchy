package api

import (
	"nikolamilovic/twitchy/accounts/api/handler"
	"nikolamilovic/twitchy/accounts/service"

	"github.com/go-playground/validator/v10"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	router         *fiber.App
	validator      *validator.Validate
	accountService service.IAccountService
}

func NewServer(service service.IAccountService) *fiber.App {
	s := &Server{
		accountService: service,
		router:         fiber.New(),
	}
	s.validator = validator.New()
	s.routes()
	return s.router
}

func (s *Server) routes() {
	h := handler.NewAuthHandler(s.validator, s.accountService)
	h.Routes()

	s.router.Mount("/api/accounts", h.Router)
}

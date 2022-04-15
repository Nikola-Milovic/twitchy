package api

import (
	"nikolamilovic/twitchy/accounts/api/handler"
	"nikolamilovic/twitchy/accounts/db"
	"nikolamilovic/twitchy/accounts/service"

	"github.com/go-playground/validator/v10"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	router         *fiber.App
	validator      *validator.Validate
	accountService service.IAccountService
	db             db.PgxIface
}

func NewServer(db db.PgxIface) *fiber.App {
	s := &Server{
		router: fiber.New(),
		db:     db,
	}
	s.validator = validator.New()
	s.routes()
	return s.router
}

func (s *Server) routes() {
	accountService := &service.AccountService{
		DB: s.db,
	}

	h := handler.NewAuthHandler(s.validator, accountService)
	h.Routes()

	s.router.Mount("/api/accounts", h.Router)
}

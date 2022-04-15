package api

import (
	"net/http"
	"nikolamilovic/twitchy/auth/api/handler"
	"nikolamilovic/twitchy/auth/db"
	"nikolamilovic/twitchy/auth/service"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	mux         *chi.Mux
	validator   *validator.Validate
	authService *service.AuthService
	db          db.PgxIface
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func NewServer(db db.PgxIface) *Server {
	s := &Server{
		mux: chi.NewMux(),
		db:  db,
	}
	s.validator = validator.New()
	s.routes()
	return s
}

func (s *Server) routes() {
	tokenService := &service.TokenService{
		DB: s.db,
	}

	authService := &service.AuthService{
		DB:           s.db,
		TokenService: tokenService,
	}

	h := handler.NewAuthHandler(s.validator, authService, tokenService)
	h.Routes()

	s.mux.Mount("/api/auth", h)
}

package api

import (
	"net/http"
	"nikolamilovic/twitchy/auth/api/handler"
	"nikolamilovic/twitchy/auth/db"
	"nikolamilovic/twitchy/auth/emitter"
	"nikolamilovic/twitchy/auth/service"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"

	ampq "github.com/rabbitmq/amqp091-go"
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

func NewServer(db db.PgxIface, ampq *ampq.Connection) (*Server, error) {
	s := &Server{
		mux: chi.NewMux(),
		db:  db,
	}
	s.validator = validator.New()

	tokenService := &service.TokenService{
		DB: s.db,
	}

	emitter, err := emitter.NewAccountEmitter(ampq)

	if err != nil {
		return nil, err
	}

	authService := &service.AuthService{
		DB:           s.db,
		TokenService: tokenService,
		Emitter:      emitter,
	}

	//Routing
	h := handler.NewAuthHandler(s.validator, authService, tokenService)
	h.Routes()

	s.mux.Mount("/v1/auth", h)
	return s, nil
}

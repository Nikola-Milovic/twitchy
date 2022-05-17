package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nikolamilovic/twitchy/auth/model/response"
	"nikolamilovic/twitchy/auth/service"
	"nikolamilovic/twitchy/common/utils"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	router       *chi.Mux
	validator    *validator.Validate
	authService  service.IAuthService
	tokenService service.ITokenService
}

func NewAuthHandler(validator *validator.Validate, auth service.IAuthService, token service.ITokenService) *AuthHandler {
	h := &AuthHandler{}

	h.authService = auth
	h.tokenService = token
	h.validator = validator

	h.Routes()

	return h
}

func (h *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *AuthHandler) Routes() {
	r := chi.NewRouter()
	h.router = r

	r.Post("/register", h.handleRegistration())
	r.Post("/login", h.handleLogin())
	r.Post("/refresh", h.handleRefresh())
}

func (h *AuthHandler) handleRegistration() http.HandlerFunc {
	type RegistrationRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
		Username string `json:"username" validate:"required"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		//https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
		var req RegistrationRequest

		if err := utils.DecodeJSONBody(w, r, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := h.validator.Struct(req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		jwt, refresh, id, err := h.authService.Register(req.Email, req.Password, req.Username)

		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := response.AuthResponse{
			JWT:          jwt,
			RefreshToken: refresh,
			ID:           id}

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			fmt.Print(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (h *AuthHandler) handleLogin() http.HandlerFunc {
	type LoginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest

		if err := utils.DecodeJSONBody(w, r, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := h.validator.Struct(req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		jwt, refresh, id, err := h.authService.Login(req.Email, req.Password)

		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := response.AuthResponse{
			JWT:          jwt,
			RefreshToken: refresh,
			ID:           id,
		}

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			fmt.Print(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (h *AuthHandler) handleRefresh() http.HandlerFunc {
	type RefreshRequest struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req RefreshRequest

		if err := utils.DecodeJSONBody(w, r, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := h.validator.Struct(req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		jwt, refresh, err := h.tokenService.RefreshToken(req.RefreshToken)

		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := response.RefreshResponse{
			JWT:          jwt,
			RefreshToken: refresh,
		}

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			fmt.Print(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

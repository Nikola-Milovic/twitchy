package handler

import (
	"net/http"
	"nikolamilovic/twitchy/accounts/service"
	"nikolamilovic/twitchy/accounts/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	Router         *fiber.App
	validator      *validator.Validate
	accountService service.IAccountService
}

func NewAuthHandler(validator *validator.Validate, accounts service.IAccountService) *AuthHandler {
	h := &AuthHandler{}

	h.accountService = accounts
	h.validator = validator

	h.Routes()

	return h
}

func (h *AuthHandler) Routes() {
	r := fiber.New()
	h.Router = r

	r.Post("/test", h.handleTest())
}

func (h *AuthHandler) handleTest() fiber.Handler {
	type TestRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
	}

	return func(ctx *fiber.Ctx) error {
		var req TestRequest

		if err := utils.DecodeJSONBody(ctx, &req); err != nil {
			return fiber.NewError(http.StatusBadRequest, err.Error())
		}

		if err := h.validator.Struct(req); err != nil {
			return fiber.NewError(http.StatusBadRequest, err.Error())
		}

		ctx.SendStatus(200)
		ctx.Send([]byte("hello there"))

		return nil
	}
}

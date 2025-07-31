package handlers

import (
	"runway/integrations"
	"runway/views/auth"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	email *email.EmailService
}

func NewAuthHandler(s *email.EmailService) *AuthHandler {
	return &AuthHandler{
		email: s,
	}
}

func (h *AuthHandler) Register(c echo.Context) error {
	if c.Request().Method == "POST" {
		input := c.FormValue("email")
		h.email.Register(input)
	}

	view := auth.Register()
	return renderView(c, auth.RegisterPage(view))
}

func (h *AuthHandler) Login(c echo.Context) error {
	view := auth.Login()
	return renderView(c, auth.LoginPage(view))
}

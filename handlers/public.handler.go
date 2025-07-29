package handlers

import (
	"github.com/labstack/echo/v4"
	"runway/views/home"
)

type PublicHandler struct{}

func NewPublicHandler() *PublicHandler {
	return &PublicHandler{}
}

func (h *PublicHandler) Index(c echo.Context) error {
	view := home.Page(home.Home())
	return renderView(c, view)

}

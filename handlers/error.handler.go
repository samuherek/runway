package handlers

import (
	"github.com/labstack/echo/v4"
	"runway/views/error_views"
)

type ErrorHnalder struct{}

func NewErrorHandler() *ErrorHnalder {
	return &ErrorHnalder{}
}

func (h *ErrorHnalder) NotFoundHandler(c echo.Context) error {
	view := error_views.Error404()
	return renderView(c, view)
}

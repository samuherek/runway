package handlers

import (
	"runway/db"
	"runway/views/app"

	"github.com/labstack/echo/v4"
)

type AppHandler struct {
	db *db.DbService
}

func NewAppHandler(db *db.DbService) *AppHandler {
	return &AppHandler{
		db: db,
	}
}

func (h *AppHandler) Home(c echo.Context) error {
	view := app.Home()
	return renderView(c, app.HomePage(view))
}

package handlers

import (
	"context"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func renderView(c echo.Context, cmp templ.Component) error {
	ctx := context.WithValue(c.Request().Context(), "H", c.Get(string("h")))
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return cmp.Render(ctx, c.Response().Writer)
}

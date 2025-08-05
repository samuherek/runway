package handlers

import (
	"context"
	"fmt"
	"net/http"
	"runway/services/notifications"

	"github.com/a-h/templ"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func renderView(c echo.Context, cmp templ.Component) error {
	ctx := context.WithValue(c.Request().Context(), "H", c.Get(string("h")))
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return cmp.Render(ctx, c.Response().Writer)
}

func unexpectedErrResponse(c echo.Context, n *notifications.Notifications) error {
	n.AddUnexpectedError()
	setHxErrors(c, n)
	return c.NoContent(http.StatusBadRequest)
}

func validationResponse(c echo.Context, n *notifications.Notifications, msgs []string) error {
	for _, m := range msgs {
		n.AddError(m)
	}
	setHxErrors(c, n)
	return c.NoContent(http.StatusBadRequest)
}

func isHxReq(c echo.Context) bool {
	return c.Request().Header.Get("HX-Request") == "true"
}

func notHxResponse(c echo.Context, n *notifications.Notifications) error {
	n.AddError("Endpoint allows only HTMX requests")
	setHxErrors(c, n)
	return c.NoContent(http.StatusNotFound)
}

func setHxErrors(c echo.Context, n *notifications.Notifications) {
	errs, _ := n.JsonErrors()
	c.Response().Header().Set(string(HxErrors), string(errs))
}

func intoValidationMessages(err error) []string {
	var msgs []string
	if verrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range verrs {
			msgs = append(msgs, fmt.Sprintf("%s: %s", e.Field(), e.Tag()))
		}
	} else {
		msgs = append(msgs, "Unexpected validation error")
	}

	return msgs
}

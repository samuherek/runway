package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"runway/handlers"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	publicH := handlers.NewPublicHandler()
	errorH := handlers.NewErrorHandler()

	e.Static("/static", "static")
	e.GET("/", publicH.Index)
	e.GET("/*", errorH.NotFoundHandler)

	e.Logger.Fatal(e.Start(":1234"))
}

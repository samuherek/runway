package main

import (
	"fmt"
	"os"
	"runway/db"
	"runway/handlers"
	email "runway/integrations"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type EnvKey string

func getEnv() EnvKey {
	env, ok := os.LookupEnv("APP_ENVIRONMENT")
	if !ok {
		fmt.Println("WARN: Did not find APP_ENVIROMENT, using development")
		return EnvKey("development")
	}
	return EnvKey(env)
}

func main() {
	e := echo.New()

	emailS := email.NewEmailService()
	dbS := db.NewDbService()
	defer dbS.Close()

	publicH := handlers.NewPublicHandler()
	authH := handlers.NewAuthHandler(emailS)
	errorH := handlers.NewErrorHandler()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/static", "static")
	e.GET("/", publicH.Index)

	e.GET("/login", authH.Login)
	e.GET("/register", authH.Register)
	e.POST("/register", authH.Register)

	e.GET("/*", errorH.NotFoundHandler)

	e.Logger.Fatal(e.Start(":1234"))
}

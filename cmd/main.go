package main

import (
	"fmt"
	"os"
	"runway/db"
	"runway/handlers"
	email "runway/integrations"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	e := echo.New()

	emailS := email.NewEmailService()
	dbS := db.NewDbService()
	defer dbS.Close()

	publicH := handlers.NewPublicHandler()
	authH := handlers.NewAuthHandler(emailS, dbS)
	errorH := handlers.NewErrorHandler()

	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/static", "static")
	e.GET("/", publicH.Index)

	e.GET("/login", authH.Login)
	e.GET("/register", authH.Register)
	e.POST("/register", authH.Register)

	e.GET("/*", errorH.NotFoundHandler)

	log.Info().Msg("Starting server")
	log.Fatal().Err(e.Start(":1234")).Msg("Server stopped")
	// e.Logger.Fatal(e.Start(":1234"))
}

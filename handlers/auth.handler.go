package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"runway/db"
	"runway/db/dbgen"
	"runway/integrations"
	"runway/utils"
	"runway/views/auth"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type AuthHandler struct {
	email *email.EmailService
	db    *db.DbService
}

func NewAuthHandler(s *email.EmailService, db *db.DbService) *AuthHandler {
	return &AuthHandler{
		email: s,
		db:    db,
	}
}

type RegisterInput struct {
	Email string `form:"email" validate:"required,email"`
}

func (h *AuthHandler) Register(c echo.Context) error {
	res := func() error {
		view := auth.Register()
		return renderView(c, auth.RegisterPage(view))
	}

	if c.Request().Method == "POST" {
		var input RegisterInput

		if err := c.Bind(&input); err != nil {
			log.Error().Err(err).Msg("Failed input binding")
			return res()
		}

		if err := c.Validate(&input); err != nil {
			// TODO: Report to UI
			fmt.Printf("Validation failed: %v\n", err)
			return res()
		}

		log.Info().Msgf("Email is now, %v", input.Email)

		_, err := h.db.Queries.GetUserByEmail(c.Request().Context(), input.Email)

		if err == nil {
			// TODO:: Message to UI that user already exists
			return res()
		}

		if !errors.Is(err, sql.ErrNoRows) {
			log.Error().Err(err).Msg("Failed db query")
			return res()
		}

		tx, err := h.db.BeginTx(c.Request().Context())
		if err != nil {
			log.Error().Err(err).Msg("Failed creating tx")
			return res()
		}

		qtx := dbgen.New(tx)

		user, err := qtx.CreateUser(c.Request().Context(), dbgen.CreateUserParams{
			ID:    uuid.New(),
			Email: input.Email,
		})

		if err != nil {
			log.Error().Err(err).Msg("Failed db create user")
			return res()
		}

		tokenValue, err := utils.GenerateToken(32)
		if err != nil {
			log.Error().Err(err).Msg("Failed token generation")
			return res()
		}

		token, err := qtx.CreateTempToken(c.Request().Context(), dbgen.CreateTempTokenParams{
			ID:        uuid.New(),
			ExpiresAt: time.Now().Add(15 * time.Minute),
			UserID: uuid.NullUUID{
				UUID:  user.ID,
				Valid: true,
			},
			Value: tokenValue,
		})

		if err != nil {
			log.Error().Err(err).Msg("Failed db create temp token")
			return res()
		}

		if err := tx.Commit(); err != nil {
			log.Error().Err(err).Msg("Failed commit transaction")
			return res()
		}

		if err := h.email.Register(input.Email, token.Value); err != nil {
			log.Error().Err(err).Msg("Failed email sending")
			return res()
		}

		// TODO: Redirect to success
	}

	return res()
}

func (h *AuthHandler) Login(c echo.Context) error {
	view := auth.Login()
	return renderView(c, auth.LoginPage(view))
}

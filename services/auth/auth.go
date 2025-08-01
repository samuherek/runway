package auth

import (
	"database/sql"
	"errors"
	"net/http"
	"runway/db"
	"runway/db/dbgen"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

const COOKIE_SESSION = "session"
const USER_ID = "user_id"

func getSessionUser(c echo.Context, db *db.DbService, token string) *uuid.UUID {
	session, err := db.Queries.GetSessionByToken(c.Request().Context(), dbgen.GetSessionByTokenParams{
		Token:     token,
		ExpiresAt: time.Now(),
	})

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Error().Err(err).Msg("Failed query db sessions")
		}
		return nil
	}

	if session.ExpiresAt.Before(time.Now()) {
		return nil
	}

	return &session.UserID
}

func AuthMiddleware(db *db.DbService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(COOKIE_SESSION)

			if err != nil || cookie.Value == "" {
				if err != nil {
					log.Error().Err(err).Msg("Failed getting cookie")
				}
				return c.Redirect(http.StatusSeeOther, "/login")
			}

			user := getSessionUser(c, db, cookie.Value)
			if user == nil {
				return c.Redirect(http.StatusSeeOther, "/login")
			}

			c.Set(USER_ID, user.ID)

			return next(c)
		}
	}
}

func RedirectIfAuthenticatedMiddleware(db *db.DbService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(COOKIE_SESSION)

			if err != nil {
				log.Error().Err(err).Msg("Failed getting cookie")
			}

			if err == nil && cookie.Value != "" {
				user := getSessionUser(c, db, cookie.Value)
				if user != nil {
					return c.Redirect(http.StatusSeeOther, "/a")
				}
			}
			return next(c)
		}
	}
}

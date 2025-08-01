package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"runway/db"
	"runway/db/dbgen"
	a "runway/services/auth"
	"runway/services/email"
	"runway/utils"
	"runway/views/auth"
	"runway/views/error_views"
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

func (h *AuthHandler) GetRegister(c echo.Context) error {
	view := auth.Register()
	return renderView(c, auth.RegisterPage(view))
}

func (h *AuthHandler) PostRegister(c echo.Context) error {
	var input RegisterInput

	res := func() error {
		view := auth.Register()
		return renderView(c, auth.RegisterPage(view))
	}

	if err := c.Bind(&input); err != nil {
		log.Error().Err(err).Msg("Failed input binding")
		return res()
	}

	if err := c.Validate(&input); err != nil {
		// TODO: Report to UI
		fmt.Printf("Validation failed: %v\n", err)
		return res()
	}

	_, err := h.db.Queries.GetUserByEmail(c.Request().Context(), input.Email)

	if err == nil {
		// TODO:: Message to UI that user already exists
		// TODO: CHeck if the user is verified. If yes, then it is an error otherwise, resend the registration link again. Probably expired
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

	return res()
}

type RegisterConfirmParams struct {
	Token string `query:"token" validate:"required,min=40,max=100`
}

func (h *AuthHandler) GetRegisterConfirm(c echo.Context) error {
	var params RegisterConfirmParams
	ctx := c.Request().Context()

	if err := c.Bind(&params); err != nil {
		log.Error().Err(err).Msg("Failed input binding")
		return renderView(c, auth.RegisterConfirmPage(auth.RegisterConfirmError("Does not look like valid link")))
	}

	if err := c.Validate(&params); err != nil {
		return renderView(c, auth.RegisterConfirmPage(auth.RegisterConfirmError("Does not look like valid link")))
	}

	token, err := h.db.Queries.GetTempToken(ctx, params.Token)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return renderView(c, auth.RegisterConfirmPage(auth.RegisterConfirmError("Does not look like valid token")))
		}
		return renderView(c, auth.RegisterConfirmPage(auth.RegisterConfirmError("Unexpected error")))
	}

	if !token.UserID.Valid {
		return renderView(c, auth.RegisterConfirmPage(auth.RegisterConfirmError("Sorry. This token does not look right.")))
	}

	if token.Used {
		return renderView(c, auth.RegisterConfirmPage(auth.RegisterConfirmError("Sorry but this token has been used already. Try again.")))
	}

	if time.Now().After(token.ExpiresAt) {
		return renderView(c, auth.RegisterConfirmPage(auth.RegisterConfirmError("Sorry but this token expired. Try again.")))
	}

	// Without transaction
	// - mark the token as used
	// if something happens here, we basically discard the token and make the user do it again
	if _, err := h.db.Queries.SetTempTokenUsed(ctx, token.Value); err != nil {
		log.Error().Err(err).Msg("Failed db setting token")
		return renderView(c, auth.RegisterConfirmPage(auth.RegisterConfirmError("Unexpected error")))
	}

	tx, err := h.db.BeginTx(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed creating tx")
		return renderView(c, auth.RegisterConfirmPage(auth.RegisterConfirmError("Unexpected error")))
	}

	qtx := dbgen.New(tx)

	qtx.SetUserVerified(ctx, dbgen.SetUserVerifiedParams{
		ID: token.UserID.UUID,
		VerifiedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})

	sessionToken, err := utils.GenerateToken(32)
	if err != nil {
		log.Error().Err(err).Msg("Failed token generation")
		return renderView(c, auth.RegisterConfirmPage(auth.RegisterConfirmError("Unexpected error")))
	}

	ua := c.Request().UserAgent()
	ip := c.RealIP()

	session, _ := qtx.CreateSession(ctx, dbgen.CreateSessionParams{
		ID:     uuid.New(),
		UserID: token.UserID.UUID,
		Token:  sessionToken,
		IpAddress: sql.NullString{
			String: ip,
			Valid:  ip != "",
		},
		UserAgent: sql.NullString{
			String: ua,
			Valid:  ua != "",
		},
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour), // 30 days
	})

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Msg("Failed commit transaction")
		return renderView(c, auth.RegisterConfirmPage(auth.RegisterConfirmError("Unexpected error")))
	}

	cookie := &http.Cookie{
		Name:     a.COOKIE_SESSION,
		Value:    session.Token,
		Path:     "/",
		Expires:  session.ExpiresAt, // 30 days
		HttpOnly: true,
		// TODO: This depends on the environment!!!!!!
		Secure:   false, // set to false if you're developing over HTTP
		SameSite: http.SameSiteLaxMode,
	}

	c.SetCookie(cookie)

	return renderView(c, auth.RegisterConfirmPage(auth.RegisterConfirm()))
}

func (h *AuthHandler) GetLogin(c echo.Context) error {
	view := auth.Login()
	return renderView(c, auth.LoginPage(view))
}

type LoginInput struct {
	Email string `form:"email" validate:"required,email"`
}

func (h *AuthHandler) PostLogin(c echo.Context) error {
	var input LoginInput
	ctx := c.Request().Context()

	res := func() error {
		view := auth.Login()
		return renderView(c, auth.LoginPage(view))
	}

	if err := c.Bind(&input); err != nil {
		log.Error().Err(err).Msg("Failed input binding")
		return res()
	}

	if err := c.Validate(&input); err != nil {
		// TODO: Report to UI
		fmt.Printf("Validation failed: %v\n", err)
		return res()
	}

	user, err := h.db.Queries.GetUserVerified(ctx, input.Email)
	if err != nil {
		// TODO: Show UI error as missing user
		return res()
	}

	sessionToken, err := utils.GenerateToken(32)
	if err != nil {
		log.Error().Err(err).Msg("Failed token generation")
		return renderView(c, auth.RegisterConfirmPage(auth.RegisterConfirmError("Unexpected error")))
	}

	ua := c.Request().UserAgent()
	ip := c.RealIP()

	session, _ := h.db.Queries.CreateSession(ctx, dbgen.CreateSessionParams{
		ID:     uuid.New(),
		UserID: user.ID,
		Token:  sessionToken,
		IpAddress: sql.NullString{
			String: ip,
			Valid:  ip != "",
		},
		UserAgent: sql.NullString{
			String: ua,
			Valid:  ua != "",
		},
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour), // 30 days
	})

	cookie := &http.Cookie{
		Name:     a.COOKIE_SESSION,
		Value:    session.Token,
		Path:     "/",
		Expires:  session.ExpiresAt, // 30 days
		HttpOnly: true,
		// TODO: This depends on the environment!!!!!!
		Secure:   false, // set to false if you're developing over HTTP
		SameSite: http.SameSiteLaxMode,
	}

	c.SetCookie(cookie)

	return c.Redirect(http.StatusSeeOther, "/a")
}

func (h *AuthHandler) GetLogout(c echo.Context) error {
	cookie, err := c.Cookie(a.COOKIE_SESSION)
	ctx := c.Request().Context()
	userId := c.Get(a.USER_ID).(*uuid.UUID)

	if err != nil {
		log.Error().Err(err).Msg("Failed getting cookie")
		return renderView(c, error_views.Error500())
	}

	// At this point, we assume we have the token in the cookie, as it passed the auth middlware
	if err := h.db.Queries.RemoveSessionByToken(ctx, dbgen.RemoveSessionByTokenParams{
		Token:  cookie.Value,
		UserID: *userId,
	}); err != nil {
		log.Error().Err(err).Msg("Failed removing session")
		return renderView(c, error_views.Error500())
	}

	c.SetCookie(&http.Cookie{
		Name:     a.COOKIE_SESSION,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
	})

	return c.Redirect(http.StatusSeeOther, "/")
}

package email

import (
	"fmt"
	l "log"
	"os"

	"github.com/resend/resend-go/v2"
	"github.com/rs/zerolog/log"
)

// TODO: Dynamic links beased on environment

func getEmailKey() string {
	env, ok := os.LookupEnv("APP_EMAIL_KEY")
	if !ok {
		l.Fatal("ERROR: Failed to load APP_EMAIL_KEY.")
	}
	return env
}

type EmailService struct {
	client *resend.Client
}

func NewEmailService() *EmailService {
	key := getEmailKey()
	client := resend.NewClient(key)

	return &EmailService{
		client,
	}
}

func (s *EmailService) Register(email, token string) error {
	params := &resend.SendEmailRequest{
		From:    "noreply@rway.app",
		To:      []string{email},
		Subject: "Let's register you",
		Html:    fmt.Sprintf("<p>Thanks for trying out RUNWAY. Before we can proceed, please register by clicking the link below!</p><p><a href='http://localhost:1234/register/confirm?token=%s'>TO BE LINK!</a></p>", token),
	}

	_, err := s.client.Emails.Send(params)

	if err != nil {
		log.Error().Err(err).Msg("Failed sending register email")
	}

	return err
}

func (s *EmailService) Login(email, token string) error {
	params := &resend.SendEmailRequest{
		From:    "noreply@rway.app",
		To:      []string{email},
		Subject: "Let's log you in",
		Html:    fmt.Sprintf("<p><p>Here is your login link ;) <a href='http://localhost:1234/login/confirm?token=%s'>TO BE LINK!</a></p>", token),
	}

	_, err := s.client.Emails.Send(params)

	if err != nil {
		log.Error().Err(err).Msg("Failed sending login email")
	}

	return err
}

package email

import (
	"fmt"
	"log"
	"os"

	"github.com/resend/resend-go/v2"
)

func getEmailKey() string {
	env, ok := os.LookupEnv("APP_EMAIL_KEY")
	if !ok {
		log.Fatal("ERROR: Failed to load APP_EMAIL_KEY.")
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
		Html:    fmt.Sprintf("<p>Thanks for trying out RUNWAY. Before we can proceed, please register by clicking the link below!</p><p><a href='http://localhost:5432/register/confirm?token=%s'>TO BE LINK!</a></p>", token),
	}

	_, err := s.client.Emails.Send(params)
	if err != nil {
		fmt.Printf("EMAIL: Register email sent\n")
	}
	return err
}

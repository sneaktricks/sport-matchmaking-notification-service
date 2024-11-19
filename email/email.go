package email

import (
	"fmt"
	"os"
	"strings"

	"github.com/sneaktricks/sport-matchmaking-notification-service/log"
	"github.com/wneessen/go-mail"
)

func NewClient() (client *mail.Client, err error) {
	if err := checkEnv(); err != nil {
		log.Logger.Error(err.Error())
	}

	client, err = mail.NewClient(
		SMTPHost,
		mail.WithPort(587),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(SMTPUsername),
		mail.WithPassword(SMTPPassword),
		mail.WithTLSPolicy(mail.TLSMandatory), // mail.WithTLSPolicy(mail.TLSOpportunistic)
	)

	return client, err
}

var (
	SMTPHost     = os.Getenv("SMTP_HOST")
	SMTPUsername = os.Getenv("SMTP_USERNAME")
	SMTPPassword = os.Getenv("SMTP_PASSWORD")
)

func checkEnv() error {
	missingEnvs := make([]string, 0)
	if SMTPHost == "" {
		missingEnvs = append(missingEnvs, "SMTP_HOST")
	}
	if SMTPUsername == "" {
		missingEnvs = append(missingEnvs, "SMTP_USERNAME")
	}
	if SMTPPassword == "" {
		missingEnvs = append(missingEnvs, "SMTP_PASSWORD")
	}

	if len(missingEnvs) > 0 {
		return fmt.Errorf("the following email environment variables are undefined: %s", strings.Join(missingEnvs, ", "))
	}

	return nil
}

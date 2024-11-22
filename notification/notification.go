package notification

import (
	"context"
	"fmt"
	"html/template"
	stdlog "log"
	"log/slog"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/sneaktricks/sport-matchmaking-notification-service/email"
	"github.com/sneaktricks/sport-matchmaking-notification-service/log"
	"github.com/sneaktricks/sport-matchmaking-notification-service/model"
	"github.com/wneessen/go-mail"
)

var matchUpdateEmailTemplate *template.Template

func init() {
	var err error
	matchUpdateEmailTemplate, err = template.ParseFiles("template/match-update-email.html.tmpl")
	if err != nil {
		stdlog.Fatalf("Failed to parse match update email template: %s", err.Error())
	}
}

type NotificationClient interface {
	SendMatchUpdateNotificationToUsers(ctx context.Context, users []*gocloak.User, details model.MatchDetails) error
}

type EmailNotificationClient struct {
	mailClient *mail.Client
}

func NewEmailNotificationClient(mc *mail.Client) *EmailNotificationClient {
	return &EmailNotificationClient{
		mailClient: mc,
	}
}

func (enc *EmailNotificationClient) SendMatchUpdateNotificationToUsers(ctx context.Context, users []*gocloak.User, details model.MatchDetails) error {
	messages := make([]*mail.Msg, len(users))

	for i, user := range users {
		msg := mail.NewMsg()
		msg.From(email.SMTPUsername)
		msg.To(*user.Email)
		msg.Subject(
			fmt.Sprintf(
				"[Sport Matchmaking] Update on your upcoming %s match (ID %s)",
				details.Sport,
				details.ID,
			),
		)
		msg.SetBodyHTMLTemplate(matchUpdateEmailTemplate, model.NewMatchUpdateTemplateData(details, user))

		messages[i] = msg
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := enc.mailClient.DialAndSendWithContext(ctx, messages...); err != nil {
		log.Logger.Error("Failed to send notification email", slog.String("error", err.Error()))
		return err
	}

	log.Logger.Info("Notification emails successfully sent", slog.Int("count", len(messages)))

	return nil
}

package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Nerzal/gocloak/v13"
	"github.com/labstack/echo/v4"
	"github.com/sneaktricks/sport-matchmaking-notification-service/auth"
	"github.com/sneaktricks/sport-matchmaking-notification-service/log"
	"github.com/sneaktricks/sport-matchmaking-notification-service/model"
)

func (h *Handler) NotifyUsersAboutMatchUpdate(c echo.Context) error {
	notificationDetails := model.NotificationDetails{}

	if err := c.Bind(&notificationDetails); err != nil {
		return HTTPError(err)
	}
	if err := c.Validate(notificationDetails); err != nil {
		return HTTPError(err)
	}

	// Get access token for retrieving users
	token, err := h.goCloakClient.GetToken(c.Request().Context(), auth.Realm, gocloak.TokenOptions{
		ClientID:     &auth.ClientID,
		ClientSecret: &auth.ClientSecret,
		GrantType:    gocloak.StringP("urn:ietf:params:oauth:grant-type:uma-ticket"),
		Audience:     gocloak.StringP(auth.ClientID),
	})
	if err != nil {
		log.Logger.Error("failed to retrieve token", slog.String("error", err.Error()))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Retrieve users with verified emails
	users, err := h.goCloakClient.GetUsers(c.Request().Context(), token.AccessToken, auth.Realm, gocloak.GetUsersParams{
		EmailVerified: gocloak.BoolP(true),
	})
	if err != nil {
		log.Logger.Error("failed to retrieve users from Keycloak", slog.String("error", err.Error()))
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	log.Logger.Debug("Found users", slog.Any("users", users))

	// Collect users to notify
	usersToNotify := make([]*gocloak.User, 0)
	for _, user := range users {
		for _, userID := range notificationDetails.UserIDs {
			if user.ID != nil && *user.ID == userID {
				usersToNotify = append(usersToNotify, user)
			}
		}
	}
	log.Logger.Debug("Users to notify", slog.Any("users", usersToNotify))

	// Send notification to users specified in NotificationDetails
	ctx := context.Background()
	go h.notificationClient.SendMatchUpdateNotificationToUsers(
		ctx,
		usersToNotify,
		notificationDetails.MatchDetails,
	)

	return c.NoContent(http.StatusNoContent)
}

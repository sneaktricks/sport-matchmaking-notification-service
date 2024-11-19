package handler

import (
	"net/http"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/labstack/echo/v4"
	"github.com/sneaktricks/sport-matchmaking-notification-service/model"
	"github.com/sneaktricks/sport-matchmaking-notification-service/notification"
)

type Handler struct {
	notificationClient notification.NotificationClient
	goCloakClient      *gocloak.GoCloak
}

func New(notificationClient notification.NotificationClient, goCloakClient *gocloak.GoCloak) *Handler {
	return &Handler{
		notificationClient: notificationClient,
		goCloakClient:      goCloakClient,
	}
}

func (h *Handler) RegisterRoutes(g *echo.Group) {
	g.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello world!")
	})

	g.GET("/time", func(c echo.Context) error {
		return c.JSON(http.StatusOK, model.TimeResponse{Time: time.Now().UTC()})
	})

	matchGroup := g.Group("/notify")
	matchGroup.GET("", h.NotifyUsersAboutMatchUpdate)
}

package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type NotificationDetails struct {
	UserIDs      []string     `json:"userIds"`
	MatchDetails MatchDetails `json:"matchDetails"`
}

type MatchDetails struct {
	ID                string    `json:"id"`
	Sport             string    `json:"sport"`
	MinParticipants   *int32    `json:"minParticipants"`
	MaxParticipants   *int32    `json:"maxParticipants"`
	StartsAt          time.Time `json:"startsAt"`
	EndsAt            time.Time `json:"endsAt"`
	Location          string    `json:"location"`
	Description       string    `json:"description"`
	ParticipationFee  int64     `json:"participationFee"`
	RequiredEquipment []string  `json:"requiredEquipment"`
	Level             string    `json:"level"`
	ChatLink          string    `json:"chatLink"`
	HostUserID        string    `json:"hostUserId"`
}

type MatchUpdateTemplateData struct {
	MatchDetails

	FormattedParticipationFee  string
	FormattedStartsAt          string
	FormattedEndsAt            string
	FormattedDuration          string
	FormattedParticipantsRange string

	Recipient *gocloak.User
}

func NewMatchUpdateTemplateData(details MatchDetails, user *gocloak.User) MatchUpdateTemplateData {
	p := message.NewPrinter(language.AmericanEnglish)

	data := MatchUpdateTemplateData{
		MatchDetails: details,
		Recipient:    user,
	}

	// Format match start and end date
	data.FormattedStartsAt = details.StartsAt.Format("Monday, January 02 at 15:04 MST")
	data.FormattedEndsAt = details.EndsAt.Format("Monday, January 02 at 15:04 MST")
	data.FormattedDuration = details.EndsAt.Sub(details.StartsAt).Truncate(1 * time.Minute).String()
	if data.FormattedDuration != "0s" {
		data.FormattedDuration = strings.TrimSuffix(data.FormattedDuration, "0s")
	}

	// Format participation fee
	fmt.Println(details)
	data.FormattedParticipationFee = "Free"
	if details.ParticipationFee > 0 {
		fmt.Println(details.ParticipationFee)
		amount := currency.EUR.Amount(float64(details.ParticipationFee) / 100)
		data.FormattedParticipationFee = p.Sprint(currency.Symbol(amount))
	}

	// Format participant range
	if data.MinParticipants != nil && data.MaxParticipants != nil {
		data.FormattedParticipantsRange = p.Sprintf("%d-%d", *data.MinParticipants, *data.MaxParticipants)
	} else if data.MinParticipants == nil && data.MaxParticipants == nil {
		data.FormattedParticipantsRange = "Any"
	} else if data.MinParticipants == nil {
		data.FormattedParticipantsRange = p.Sprintf("≤%d", *details.MaxParticipants)
	} else {
		data.FormattedParticipantsRange = p.Sprintf("%d+", *details.MinParticipants)
	}

	return data
}

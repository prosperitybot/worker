package internal

import (
	"encoding/json"
	"net/http"

	"github.com/prosperitybot/worker/model"
)

type ResponseType uint64

const (
	AccessDenied ResponseType = iota
	PremiumOnly
	Successful
)

func RespondToRequest(w http.ResponseWriter, status int, body interface{}) error {
	// Writes the status code passed through as the header for the response.
	w.WriteHeader(status)

	// Generates a json object based on the interface that is passed in and also specifies the status code in the response.
	return json.NewEncoder(w).Encode(body)
}

func RespondToInteraction(w http.ResponseWriter, status ResponseType, message string) error {
	var color uint32

	switch status {
	case AccessDenied:
		color = 16711680
		break
	case PremiumOnly:
		color = 10494192
		break
	default:
		color = 37119
		break
	}

	return json.NewEncoder(w).Encode(model.InteractionResponse{
		Type: model.ChannelMessageWithSourceCallback,
		Data: model.InteractionCallbackData{
			Embeds: []model.MessageEmbed{{
				Title:       "Prosperity Bot",
				Type:        model.EmbedRichType,
				Color:       color,
				Description: message,
			}},
		},
	})
}

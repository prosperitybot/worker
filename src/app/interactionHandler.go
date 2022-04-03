package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/prosperitybot/worker/internal"
	"github.com/prosperitybot/worker/model"
)

func BaseRequest(w http.ResponseWriter, r *http.Request) {
	apiRequest := model.Interaction{}

	err := json.NewDecoder(r.Body).Decode(&apiRequest)

	if err != nil {
		fmt.Printf("Error whilst decoding API request: %v", err)
		internal.RespondToRequest(w, http.StatusBadRequest, "Invalid API Request")
		return
	}

	w.Header().Set("Authorization", fmt.Sprintf("Bot %s", os.Getenv("BOT_TOKEN")))
	w.Header().Set("User-Agent", os.Getenv("BOT_USER_AGENT"))
	switch apiRequest.Type {
	case model.InteractionPing:
		internal.RespondToRequest(w, http.StatusOK, map[string]int{"type": 1})
		break
	case model.InteractionApplicationCommand:
		internal.RespondToRequest(w, http.StatusOK, model.InteractionResponse{
			Type: model.ChannelMessageWithSourceCallback,
			Data: model.InteractionCallbackData{
				Embeds: []model.MessageEmbed{{
					Title:       "Test",
					Type:        model.EmbedRichType,
					Description: "This is a test command (This bot is in test mode)",
				}},
			},
		})
	}

	commandId := apiRequest.Data.ID

	fmt.Println(commandId)
	internal.RespondToRequest(w, http.StatusOK, "commandId")
}

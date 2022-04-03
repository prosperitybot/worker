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

	if apiRequest.Type == 1 {
		w.Header().Set("Authorization", fmt.Sprintf("Bot %s", os.Getenv("BOT_TOKEN")))
		w.Header().Set("User-Agent", os.Getenv("USER_AGENT"))
		internal.RespondToRequest(w, http.StatusOK, map[string]int{"type": 1})
		return
	}

	commandId := apiRequest.Data.ID

	fmt.Println(commandId)
	internal.RespondToRequest(w, http.StatusOK, "commandId")
}

package app

import (
	"net/http"

	"github.com/prosperitybot/worker/internal"
)

func CheckHealth(w http.ResponseWriter, r *http.Request) {
	// result := internal.Database.Ping()

	// if result != nil {
	// 	internal.RespondToRequest(w, http.StatusInternalServerError, healthResponse{Status: http.StatusInternalServerError})
	// } else {
	// 	internal.RespondToRequest(w, http.StatusOK, healthResponse{Status: http.StatusOK})
	// }
	internal.RespondToRequest(w, http.StatusOK, healthResponse{Status: http.StatusOK})
}

type healthResponse struct {
	Status int `json:"status"`
}

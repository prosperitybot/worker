package app

import (
	"net/http"

	"github.com/prosperitybot/worker/internal"
	"github.com/prosperitybot/worker/services"
)

type TestHandler struct {
	service services.TestService
}

func (h *TestHandler) Test(w http.ResponseWriter, r *http.Request) {
	internal.RespondToRequest(w, 200, "Test")
}

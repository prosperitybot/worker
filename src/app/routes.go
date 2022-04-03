package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

func HandleRequests(router *mux.Router, testHandler TestHandler) {
	// r := router.PathPrefix("/").Subrouter()

	router.HandleFunc("/interactions", BaseRequest).Methods(http.MethodPost)

	router.HandleFunc("/health", CheckHealth).Methods(http.MethodGet)
}

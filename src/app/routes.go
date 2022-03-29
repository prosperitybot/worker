package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

func HandleRequests(router *mux.Router, testHandler TestHandler) {
	r := router.PathPrefix("/").Subrouter()
	testRouter := r.PathPrefix("/test").Subrouter()

	router.HandleFunc("/health", CheckHealth).Methods(http.MethodGet)

	testRouter.HandleFunc("", testHandler.Test).Methods(http.MethodGet)
}

package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Sets the Content-Type for the response
			w.Header().Add("Content-Type", "application/json")

			// Passes on the HTTP request to actually run
			next.ServeHTTP(w, r)
		})
	}
}

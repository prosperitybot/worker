package app

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/prosperitybot/worker/internal"
)

func NewMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Sets the Content-Type for the response
			w.Header().Add("Content-Type", "application/json")

			// Validates the payload
			if ValidatePayload(w, r) == false {
				internal.RespondToRequest(w, http.StatusUnauthorized, "invalid request signature")
				return
			}

			// Passes on the HTTP request to actually run
			next.ServeHTTP(w, r)
		})
	}
}

func ValidatePayload(w http.ResponseWriter, r *http.Request) bool {
	key, _ := hex.DecodeString(os.Getenv("BOT_PUBLIC_KEY"))
	signature := r.Header.Get("X-Signature-Ed25519")
	timestamp := r.Header.Get("X-Signature-Timestamp")
	sig, err := hex.DecodeString(signature)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	body := new(bytes.Buffer)
	if _, err := io.Copy(body, r.Body); err != nil {
		fmt.Println(err.Error())
		return false
	}
	buf := bytes.NewBufferString(timestamp + body.String())
	if !ed25519.Verify(key, buf.Bytes(), sig) {
		fmt.Println("Invalid Signature")
		return false
	}
	r.Body = io.NopCloser(body)

	return false
}

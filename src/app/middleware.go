package app

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/prosperitybot/worker/internal"
	"github.com/prosperitybot/worker/logging"
	"go.uber.org/zap"
)

func NewMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()

			// Sets the Content-Type for the response
			w.Header().Add("Content-Type", "application/json")

			logging.Info(r.Context(), fmt.Sprintf("Req: %s %s", r.Method, r.RequestURI))
			// Validates the payload
			if r.RequestURI != "/health" {
				if ValidatePayload(w, r) == false {
					internal.RespondToRequest(w, http.StatusUnauthorized, "invalid request signature")
					return
				}
			}

			logRespWriter := logging.NewLogResponseWriter(w)

			// Passes on the HTTP request to actually run
			next.ServeHTTP(logRespWriter, r)

			logging.Info(r.Context(), fmt.Sprintf("Res: %s %s", r.Method, r.RequestURI),
				zap.Int("status", logRespWriter.StatusCode),
				zap.String("response_time", fmt.Sprintf("%d ms", time.Since(startTime).Milliseconds())),
				zap.Int("response_bytes", logRespWriter.Buf.Len()),
				zap.String("method", r.Method),
				zap.String("user_agent", r.Header.Get("User-Agent")))
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
	fmt.Println(body.String())
	buf := bytes.NewBufferString(timestamp + body.String())
	if !ed25519.Verify(key, buf.Bytes(), sig) {
		fmt.Println("Invalid Signature")
		return false
	}
	r.Body = io.NopCloser(body)

	return false
}

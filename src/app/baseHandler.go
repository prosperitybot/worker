package app

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/prosperitybot/worker/internal"
)

func BaseRequest(w http.ResponseWriter, r *http.Request) {
	key, _ := hex.DecodeString(os.Getenv("BOT_PUBLIC_KEY"))
	signature := r.Header.Get("X-Signature-Ed25519")
	timestamp := r.Header.Get("X-Signature-Timestamp")
	sig, err := hex.DecodeString(signature)
	if err != nil {
		fmt.Println(err.Error())
		internal.RespondToRequest(w, http.StatusUnauthorized, err.Error())
		return
	}
	body := new(bytes.Buffer)
	if _, err := io.Copy(body, r.Body); err != nil {
		fmt.Println(err.Error())
		internal.RespondToRequest(w, http.StatusUnauthorized, err.Error())
		return
	}
	buf := bytes.NewBufferString(timestamp + body.String())
	if !ed25519.Verify(key, buf.Bytes(), sig) {
		fmt.Println("Invalid Signature")
		internal.RespondToRequest(w, http.StatusUnauthorized, "Invalid Signature")
		return
	}
	r.Body = io.NopCloser(body)
	internal.RespondToRequest(w, http.StatusOK, "Status OK")
}

package cursor

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"io"
	"net/http"
	"strings"
)

// VerifySignature validates the HMAC-SHA256 signature header against the raw body and secret.
func VerifySignature(secret string, body []byte, signatureHeader string) bool {
	const prefix = "sha256="
	if !strings.HasPrefix(signatureHeader, prefix) {
		return false
	}
	expected := signatureHeader[len(prefix):]
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	computed := hex.EncodeToString(mac.Sum(nil))
	return subtle.ConstantTimeCompare([]byte(expected), []byte(computed)) == 1
}

// Middleware returns an http.Handler that verifies webhook signatures.
// If verification fails, it responds with 401. On success, it calls next.
func Middleware(secret string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read the raw body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusBadRequest)
			return
		}
		// Restore body for downstream handlers
		r.Body = io.NopCloser(bytes.NewReader(body))

		sig := r.Header.Get("X-Webhook-Signature")
		if !VerifySignature(secret, body, sig) {
			http.Error(w, "invalid signature", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

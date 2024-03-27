package endpoint

import (
	"errors"
	"net/http"
)

var errRefreshTokenHashRequired = errors.New("refresh_token cookie not found")
var errFingerprintRequired = errors.New("fingerprint cookie not found")

func (e *Endpoint) RefreshKeyPair(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed. Only POST is allowed", http.StatusMethodNotAllowed)
		return
	}
	refreshTokenHash, fingerprint, err := extractTokenCookies(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tokens, err := e.service.RefreshKeyPair(refreshTokenHash, r.RemoteAddr, fingerprint)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &tokens.AccessToken)
	http.SetCookie(w, &tokens.RefreshToken)
	http.SetCookie(w, &tokens.Fingerprint)
	w.WriteHeader(http.StatusOK)
}

func extractTokenCookies(r *http.Request) (string, string, error) {
	var RefreshTokenHash string
	var fingerprint string
	for _, cookie := range r.Cookies() {
		if cookie.Name == "refresh_token" {
			RefreshTokenHash = cookie.Value
		} else if cookie.Name == "fingerprint" {
			fingerprint = cookie.Value
		}
	}
	if RefreshTokenHash == "" {
		return "", "", errRefreshTokenHashRequired
	} else if fingerprint == "" {
		return "", "", errFingerprintRequired
	}
	return RefreshTokenHash, fingerprint, nil
}

package endpoint

import (
	"net/http"
)

func (e *Endpoint) GetKeyPair(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed. Only GET is allowed", http.StatusMethodNotAllowed)
		return
	}
	GUID := r.URL.Query().Get("GUID")
	remoteAddr := r.RemoteAddr
	tokens, err := e.service.GetKeyPair(GUID, remoteAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &tokens.AccessToken)
	http.SetCookie(w, &tokens.RefreshToken)
	http.SetCookie(w, &tokens.Fingerprint)
	w.WriteHeader(http.StatusOK)
}

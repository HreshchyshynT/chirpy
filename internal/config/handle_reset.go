package config

import "net/http"

func (ac *ApiConfig) HandleReset(w http.ResponseWriter, r *http.Request) {
	if ac.Platform != Dev {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	ac.FileserverHits.Store(0)
	ac.Queries.ClearUsers(r.Context())
	w.WriteHeader(http.StatusOK)
}

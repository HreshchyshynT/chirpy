package config

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/hreshchyshynt/chirpy/internal/auth"
	"github.com/hreshchyshynt/chirpy/internal/utils"
)

func (ac *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ac.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (ac *ApiConfig) MiddlewareDbAccess(next DbHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next(w, r, ac.Queries)
	})
}

func (ac *ApiConfig) MiddlewareWithConfig(next ApiHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next(w, r, ac)
	})
}

func (ac *ApiConfig) MiddlewareRequireJWT(
	next func(w http.ResponseWriter, r *http.Request, userId uuid.UUID),
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "No token provided", err)
			return
		}

		id, err := auth.ValidateJWT(token, ac.Secret)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token", err)
			return
		}

		next(w, r, id)
	})
}

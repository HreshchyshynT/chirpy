package user

import (
	"net/http"
	"time"

	"github.com/hreshchyshynt/chirpy/internal/auth"
	"github.com/hreshchyshynt/chirpy/internal/config"
	"github.com/hreshchyshynt/chirpy/internal/utils"
)

func HandleRefreshToken(
	w http.ResponseWriter,
	r *http.Request,
	config *config.ApiConfig,
) {

	type responseBody struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		utils.RespondWithError(
			w,
			http.StatusBadRequest,
			"Invalid refresh token",
			err,
		)
		return
	}

	rt, err := config.Queries.FindRefreshToken(r.Context(), token)
	if err != nil || rt.RevokedAt.Valid || time.Now().After(rt.ExpiresAt) {
		utils.RespondWithError(
			w,
			http.StatusUnauthorized,
			"Invalid refresh token",
			err,
		)
		return
	}

	jwt, err := auth.MakeJWT(rt.UserID, config.Secret, time.Hour)
	if err != nil {
		utils.RespondWithError(
			w,
			http.StatusInternalServerError,
			"Can not generate JWT",
			err,
		)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, responseBody{Token: jwt})

}

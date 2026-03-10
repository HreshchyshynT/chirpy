package user

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/hreshchyshynt/chirpy/internal/auth"
	"github.com/hreshchyshynt/chirpy/internal/config"
	"github.com/hreshchyshynt/chirpy/internal/database"
	"github.com/hreshchyshynt/chirpy/internal/utils"
)

func HandleLogin(
	w http.ResponseWriter,
	r *http.Request,
	config *config.ApiConfig,
) {
	type requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type responseBody struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	var body requestBody

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	defer r.Body.Close()

	err := decoder.Decode(&body)
	if err != nil {
		utils.RespondWithError(
			w,
			http.StatusBadRequest,
			utils.MessageInvalidRequestBody,
			err,
		)
		return
	}

	user, err := config.Queries.FindUser(r.Context(), body.Email)
	if err != nil {
		utils.RespondWithError(
			w,
			http.StatusUnauthorized,
			"Invalid email or password",
			err,
		)
		return
	}

	match, err := auth.CheckPasswordHash(body.Password, user.HashedPassword)

	if !match {
		utils.RespondWithError(
			w,
			http.StatusUnauthorized,
			"Invalid email or password",
			err,
		)
		return
	}

	token, err := auth.MakeJWT(
		user.ID,
		config.Secret,
		time.Hour,
	)

	if err != nil {
		utils.RespondWithError(
			w,
			http.StatusInternalServerError,
			"Can not generate JWT",
			err,
		)
		return
	}

	refreshToken := auth.MakeRefreshToken()
	_, err = config.Queries.SaveRefreshToken(
		r.Context(),
		database.SaveRefreshTokenParams{
			Token:     refreshToken,
			UserID:    user.ID,
			ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
		},
	)

	if err != nil {
		utils.RespondWithError(
			w,
			http.StatusInternalServerError,
			"Can not generate refresh token",
			err,
		)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, responseBody{
		User:         toDomainUser(user),
		Token:        token,
		RefreshToken: refreshToken,
	})
}

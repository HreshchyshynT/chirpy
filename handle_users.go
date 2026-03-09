package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hreshchyshynt/chirpy/internal/auth"
	"github.com/hreshchyshynt/chirpy/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func toDomainUser(u database.User) User {
	return User{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Email:     u.Email,
	}
}

func handleCreateUser(
	w http.ResponseWriter,
	r *http.Request,
	db *database.Queries,
) {
	type requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body requestBody

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	defer r.Body.Close()

	err := decoder.Decode(&body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, messageInvalidRequestBody, err)
		return
	}

	pass, err := auth.HashPassword(body.Password)
	if err != nil {
		respondWithError(w,
			http.StatusInternalServerError,
			messageInternalServerError,
			err,
		)
		return
	}

	user, err := db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          body.Email,
		HashedPassword: pass,
	})

	if IsDuplicatedKeys(err) {
		respondWithError(w, http.StatusBadRequest, "Email is already used", err)
		return
	}

	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Can't create user. Try again later",
			err,
		)
		return
	}

	respondWithJSON(w, http.StatusCreated, toDomainUser(user))
}

func handleLogin(
	w http.ResponseWriter,
	r *http.Request,
	config *ApiConfig,
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
		respondWithError(w, http.StatusBadRequest, messageInvalidRequestBody, err)
		return
	}

	user, err := config.Queries.FindUser(r.Context(), body.Email)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Invalid email or password",
			err,
		)
		return
	}

	match, err := auth.CheckPasswordHash(body.Password, user.HashedPassword)

	if !match {
		respondWithError(
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
		respondWithError(
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
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Can not generate refresh token",
			err,
		)
		return
	}

	respondWithJSON(w, http.StatusOK, responseBody{
		User:         toDomainUser(user),
		Token:        token,
		RefreshToken: refreshToken,
	})
}

func handleRefreshToken(
	w http.ResponseWriter,
	r *http.Request,
	config *ApiConfig,
) {

	type responseBody struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			"Invalid refresh token",
			err,
		)
		return
	}

	rt, err := config.Queries.FindRefreshToken(r.Context(), token)
	if err != nil || rt.RevokedAt.Valid || rt.ExpiresAt.After(time.Now()) {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Invalid refresh token",
			err,
		)
		return
	}

	jwt, err := auth.MakeJWT(rt.UserID, config.Secret, time.Hour)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Can not generate JWT",
			err,
		)
		return
	}

	respondWithJSON(w, http.StatusOK, responseBody{Token: jwt})

}

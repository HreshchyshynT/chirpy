package user

import (
	"encoding/json"
	"net/http"

	"github.com/hreshchyshynt/chirpy/internal/auth"
	"github.com/hreshchyshynt/chirpy/internal/database"
	"github.com/hreshchyshynt/chirpy/internal/utils"
)

func HandleCreateUser(
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
		utils.RespondWithError(
			w,
			http.StatusBadRequest,
			utils.MessageInvalidRequestBody,
			err,
		)
		return
	}

	pass, err := auth.HashPassword(body.Password)
	if err != nil {
		utils.RespondWithError(w,
			http.StatusInternalServerError,
			utils.MessageInternalServerError,
			err,
		)
		return
	}

	user, err := db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          body.Email,
		HashedPassword: pass,
	})

	if utils.IsDuplicatedKeys(err) {
		utils.RespondWithError(w, http.StatusBadRequest, "Email is already used", err)
		return
	}

	if err != nil {
		utils.RespondWithError(
			w,
			http.StatusInternalServerError,
			"Can't create user. Try again later",
			err,
		)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, toDomainUser(user))
}

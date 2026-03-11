package user

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/hreshchyshynt/chirpy/internal/config"
	"github.com/hreshchyshynt/chirpy/internal/database"
	"github.com/hreshchyshynt/chirpy/internal/utils"
)

const (
	eventUpgraded = "user.upgraded"
)

func HandlePolkaWebhook(
	w http.ResponseWriter,
	r *http.Request,
	apiConfig *config.ApiConfig,
) {

	type RequestBody struct {
		Event string `json:"event"`
		Data  struct {
			UserId uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var request RequestBody
	err := decoder.Decode(&request)
	if err != nil {
		utils.RespondWithError(
			w,
			http.StatusBadRequest,
			utils.MessageBadRequest,
			err,
		)
		return
	}

	if request.Event != eventUpgraded {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	user, err := apiConfig.Queries.SetIsChirpyRedUser(
		r.Context(),
		database.SetIsChirpyRedUserParams{
			IsChirpyRed: true,
			UserID:      request.Data.UserId,
		},
	)

	// if chirpy is false after update - user with such id not exists
	if err != nil || !user.IsChirpyRed {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

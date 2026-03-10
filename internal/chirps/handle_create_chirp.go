package chirps

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hreshchyshynt/chirpy/internal/database"
	"github.com/hreshchyshynt/chirpy/internal/utils"
)

func HandleCreateChirp(
	w http.ResponseWriter,
	r *http.Request,
	userId uuid.UUID,
	db *database.Queries,
) {
	type requestBody struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	defer r.Body.Close()

	var request requestBody

	err := decoder.Decode(&request)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, utils.MessageInvalidRequestBody, err)
		return
	}

	cleanedText, err := validateChirp(request.Body)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	chirp, err := db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedText,
		UserID: userId,
	})
	if err != nil {
		utils.RespondWithError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Can't save chirp: %v\n", err),
			err,
		)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, toDomainChirp(chirp))

}

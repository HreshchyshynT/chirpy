package chirps

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/hreshchyshynt/chirpy/internal/database"
	"github.com/hreshchyshynt/chirpy/internal/utils"
)

func HandleGetAllChirps(
	w http.ResponseWriter,
	r *http.Request,
	db *database.Queries,
) {

	authorIdParam := r.URL.Query().Get("author_id")

	authorId, err := uuid.Parse(authorIdParam)
	if err != nil {
		authorId = uuid.Nil
	}

	chirps, err := db.GetAllChirps(r.Context(), uuid.NullUUID{
		UUID:  authorId,
		Valid: err == nil,
	})
	if err != nil {
		utils.RespondWithError(
			w,
			http.StatusInternalServerError,
			"Can't get chirps",
			err,
		)
		return
	}

	responseChirps := make([]Chirp, len(chirps))

	for i, c := range chirps {
		responseChirps[i] = toDomainChirp(c)
	}

	utils.RespondWithJSON(w, http.StatusOK, responseChirps)
}

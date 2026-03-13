package chirps

import (
	"net/http"
	"sort"

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
	sortParam := r.URL.Query().Get("sort")

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

	sort.Slice(responseChirps, func(i, j int) bool {
		left, right := responseChirps[i], responseChirps[j]
		switch sortParam {
		case "desc":
			return left.CreatedAt.After(right.CreatedAt)
		default:
			return left.CreatedAt.Before(right.CreatedAt)
		}
	})

	utils.RespondWithJSON(w, http.StatusOK, responseChirps)
}

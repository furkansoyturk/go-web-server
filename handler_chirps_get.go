package main

import (
	"net/http"
	"sort"
	"strconv"
)

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}
	authorIdParam := r.URL.Query().Get("author_id")
	authorId, _ := strconv.Atoi(authorIdParam)

	queriedChirps := []Chirp{}
	if 0 != authorId {
		for _, dbChirp := range dbChirps {
			if dbChirp.AuthorID == authorId {
				queriedChirps = append(queriedChirps, Chirp{
					ID:       dbChirp.ID,
					Body:     dbChirp.Body,
					AuthorID: dbChirp.AuthorID,
				})
			}
		}
		if len(queriedChirps) == 0 {
			respondWithError(w, http.StatusNotFound, "Author not found")
		}

		sort.Slice(queriedChirps, func(i, j int) bool {
			return queriedChirps[i].ID < queriedChirps[j].ID
		})

		respondWithJSON(w, http.StatusOK, queriedChirps)
		return
	} else {
		chirps := []Chirp{}
		for _, dbChirp := range dbChirps {
			chirps = append(chirps, Chirp{
				ID:       dbChirp.ID,
				Body:     dbChirp.Body,
				AuthorID: dbChirp.AuthorID,
			})
		}

		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID < chirps[j].ID
		})

		respondWithJSON(w, http.StatusOK, chirps)
	}
}

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := strconv.Atoi(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	dbChirp, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp")
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:       dbChirp.ID,
		Body:     dbChirp.Body,
		AuthorID: dbChirp.AuthorID,
	})
}

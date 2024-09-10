package main

import (
	"net/http"
	"strconv"

	"github.com/furkansoyturk/go-web-server/internal/auth"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := strconv.Atoi(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	subject, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Couldn't validate JWT")
		return
	}

	dbChirp, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp")
		return
	}

	author, _ := strconv.Atoi(subject)
	if dbChirp.AuthorID != author {
		respondWithError(w, http.StatusForbidden, "Unauthorized")
		return
	}

	cfg.DB.DeleteChirp(author)
	respondWithJSON(w, http.StatusNoContent, nil)
}

package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/furkansoyturk/go-web-server/internal/auth"
	"github.com/furkansoyturk/go-web-server/internal/database"
)

type UsersUpdate struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		UsersUpdate
	}
	token := r.Header.Get("Authorization")

	if len(token) > 0 {
		tokenArr := strings.Split(token, " ")
		if tokenArr[0] == "Bearer" {
			token = tokenArr[1]
		} else {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
	} else {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	id, err := strconv.Atoi(auth.ReadFrom(token))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	usr, err := cfg.DB.UpdateUser(params.Email, hashedPassword, id)
	if err != nil {
		if errors.Is(err, database.ErrAlreadyExists) {
			respondWithError(w, http.StatusConflict, "User already exists")
			return
		}

		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		UsersUpdate: UsersUpdate{
			ID:    usr.ID,
			Email: usr.Email,
		},
	})
}

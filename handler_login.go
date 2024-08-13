package main

import (
	"encoding/json"
	"net/http"
)

type LoginResponse struct {
	ID    int    `json:"id"`
	EMAIL string `json:"email"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type UserCreateRequest struct {
		EMAIL    string `json:"email"`
		PASSWORD string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := UserCreateRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.Login(params.EMAIL, params.PASSWORD)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't login")
		return
	}

	respondWithJSON(w, http.StatusOK, LoginResponse{
		ID:    user.ID,
		EMAIL: user.EMAIL,
	})
}

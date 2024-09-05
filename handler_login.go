package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/furkansoyturk/go-web-server/internal/auth"
)

type loginResponse struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	JWTToken     string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password        string `json:"password"`
		Email           string `json:"email"`
		ExpiresInSecond int    `json:"expires_in_seconds"`
	}
	type response struct {
		loginResponse
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	token, refreshToken, err := auth.CreateJWT([]byte(cfg.jwtSecret), strconv.Itoa(user.ID), params.ExpiresInSecond)
	user, err = cfg.DB.SaveRefreshToken(user.ID, refreshToken)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		loginResponse: loginResponse{
			ID:           user.ID,
			Email:        user.Email,
			JWTToken:     token,
			RefreshToken: user.RefreshToken,
		},
	})
}

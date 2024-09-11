package main

import (
	"encoding/json"
	"net/http"

	"github.com/furkansoyturk/go-web-server/internal/auth"
)

const userUpgradedEvent = "user.upgraded"

func (cfg *apiConfig) handlerWebhooks(w http.ResponseWriter, r *http.Request) {
	// TODO: fix webhooks.
	type Data struct {
		UserID int `json:"user_id"`
	}
	type parameters struct {
		Event string `json:"event"`
		Data  `json:"data"`
	}

	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil || apiKey != cfg.webhookSecret {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.GetUser(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	if params.Event == userUpgradedEvent {
		cfg.DB.UpdateUserMembership(user.ID, true)
	}

	respondWithJSON(w, http.StatusNoContent, "")
}

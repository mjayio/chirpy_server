package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/mjayio/server/internal/auth"
)

func (cfg *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	type webhook struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		}
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}
	if apiKey != cfg.apiKey {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := webhook{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		respondWithError(w, http.StatusNoContent, "Invalid event", nil)
		return
	}

	err = cfg.database.MakeChirpyRed(r.Context(), params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't make user chirpy red", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

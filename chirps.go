package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"github.com/google/uuid"
	"github.com/mjayio/server/internal/auth"
	"github.com/mjayio/server/internal/database"
)

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		ID        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Body      string `json:"body"`
		UserID    string `json:"user_id"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	wordsToReplace := []string{"kerfuffle", "sharbert", "fornax"}
	for _, word := range wordsToReplace {
		pattern := `(?i)\b` + regexp.QuoteMeta(word) + `\b`
		regex := regexp.MustCompile(pattern)
		params.Body = regex.ReplaceAllString(params.Body, "****")
	}

	chirp, err := cfg.database.CreateChirp(r.Context(), database.CreateChirpParams{
		UserID: userID,
		Body:   params.Body,
	})

	if err != nil {
		log.Printf("Error creating chirp: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, returnVals{
		ID:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt.String(),
		UpdatedAt: chirp.UpdatedAt.String(),
		Body:      chirp.Body,
		UserID:    chirp.UserID.String(),
	})
}

func (cfg *apiConfig) handlerChirpsListAuthor(w http.ResponseWriter, r *http.Request, authorID uuid.UUID, sort string) {
	chirps, err := cfg.database.ListChirpsByAuthor(r.Context(), authorID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't list chirps", err)
		return
	}

	type chirpResponse struct {
		ID        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Body      string `json:"body"`
		UserID    string `json:"user_id"`
	}

	response := make([]chirpResponse, len(chirps))

	for i, chirp := range chirps {
		response[i] = chirpResponse{
			ID:        chirp.ID.String(),
			CreatedAt: chirp.CreatedAt.String(),
			UpdatedAt: chirp.UpdatedAt.String(),
			Body:      chirp.Body,
			UserID:    chirp.UserID.String(),
		}
	}

	if sort == "desc" {
		for i, j := 0, len(response)-1; i < j; i, j = i+1, j-1 {
			response[i], response[j] = response[j], response[i]
		}
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) handlerChirpsList(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort")
	if sort != "asc" && sort != "desc" {
		respondWithError(w, http.StatusBadRequest, "Invalid sort parameter", nil)
		return
	}

	authorID := r.URL.Query().Get("author_id")
	if authorID != "" {
		parseAuthorID, err := uuid.Parse(authorID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
			return
		}
		cfg.handlerChirpsListAuthor(w, r, parseAuthorID, sort)
		return
	}

	chirps, err := cfg.database.ListChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't list chirps", err)
		return
	}

	type chirpResponse struct {
		ID        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Body      string `json:"body"`
		UserID    string `json:"user_id"`
	}

	response := make([]chirpResponse, len(chirps))
	for i, chirp := range chirps {
		response[i] = chirpResponse{
			ID:        chirp.ID.String(),
			CreatedAt: chirp.CreatedAt.String(),
			UpdatedAt: chirp.UpdatedAt.String(),
			Body:      chirp.Body,
			UserID:    chirp.UserID.String(),
		}
	}

	if sort == "desc" {
		for i, j := 0, len(response)-1; i < j; i, j = i+1, j-1 {
			response[i], response[j] = response[j], response[i]
		}
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) handlerChirpsRead(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	if chirpID == "" {
		respondWithError(w, http.StatusBadRequest, "chirpID is required", nil)
		return
	}

	parseChirpID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	chirp, err := cfg.database.GetChirp(r.Context(), parseChirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	type chirpResponse struct {
		ID        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Body      string `json:"body"`
		UserID    string `json:"user_id"`
	}

	response := chirpResponse{
		ID:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt.String(),
		UpdatedAt: chirp.UpdatedAt.String(),
		Body:      chirp.Body,
		UserID:    chirp.UserID.String(),
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	if chirpID == "" {
		respondWithError(w, http.StatusBadRequest, "chirpID is required", nil)
		return
	}

	parseChirpID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	chirp, err := cfg.database.GetChirp(r.Context(), parseChirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	if chirp.UserID.String() != userID.String() {
		respondWithError(w, http.StatusForbidden, "You can only delete your own chirps", nil)
		return
	}

	err = cfg.database.DeleteChirp(r.Context(), parseChirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

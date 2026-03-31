package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/cfranklin121/chirpy/internal/auth"
	"github.com/cfranklin121/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type RequestBody struct {
		Body string `json:"body"`
	}

	type Params struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserId    uuid.UUID `json:"user_id"`
	}
	log.Printf("%s %s", r.Method, r.URL.Path)

	decoder := json.NewDecoder(r.Body)

	reqBody := RequestBody{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Printf("%s %s", r.Method, r.URL.Path)
		respondWithError(w, 500, "Could not decode JSON")
		return
	}
	log.Println(reqBody)

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	user_id, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	refresh_token, err := cfg.db.GetRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, 401, "Invalid token")
		return
	}

	if refresh_token.RevokedAt.Valid {
		respondWithError(w, 401, "Token has been revoked")
		return
	}

	if len(reqBody.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	cleaned := cleanString(reqBody.Body)
	reqBody.Body = cleaned

	params := database.CreateChirpParams{
		Body:   reqBody.Body,
		UserID: user_id,
	}
	chirp, err := cfg.db.CreateChirp(r.Context(), params)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	respondWithJSON(w, 201, Params{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	})
}

func cleanString(str string) string {
	slice := strings.Split(str, " ")
	newStr := []string{}
	for _, word := range slice {
		switch strings.ToLower(word) {
		case "kerfuffle":
			word = "****"
		case "sharbert":
			word = "****"
		case "fornax":
			word = "****"
		}
		newStr = append(newStr, word)
	}
	return strings.Join(newStr, " ")
}

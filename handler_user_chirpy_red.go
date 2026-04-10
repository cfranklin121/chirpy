package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/cfranklin121/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUpgradeChirpyRed(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)
	defer r.Body.Close()
	type RequestBody struct {
		Event string `json:"event"`
		Data  struct {
			UserId uuid.UUID `json:"user_id"`
		} `json:"data"`
	}
	type ReturnVal struct{}

	decoder := json.NewDecoder(r.Body)
	reqBody := RequestBody{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		respondWithError(w, 500, "Could not decode")
		return
	}

	if reqBody.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, ReturnVal{})
		return
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid API Key")
		return
	}

	_, err = cfg.db.UpgradeToChirpyRed(r.Context(), reqBody.Data.UserId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No user found")
		return
	}

	respondWithJSON(w, http.StatusNoContent, ReturnVal{})
}

package main

import (
	"log"
	"net/http"

	"github.com/cfranklin121/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)
	defer r.Body.Close()

	type response struct{}

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid Token")
		return
	}

	chirpIdString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIdString)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)

	if userID != chirp.UserID {
		respondWithError(w, http.StatusForbidden, "Not authorized")
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}

	respondWithJSON(w, http.StatusNoContent, response{})
}

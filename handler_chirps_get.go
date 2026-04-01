package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	arr := []Chirp{}
	for _, chirp := range chirps {

		arr = append(arr, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID,
		})

	}
	log.Printf("%s %s", r.Method, r.URL.Path)
	respondWithJSON(w, 200, arr)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)
	chirpId := r.PathValue("chirpID")

	id, err := uuid.Parse(chirpId)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), uuid.UUID(id))
	if err != nil {
		respondWithError(w, 404, err.Error())
		return
	}

	log.Printf("%s %s", r.Method, r.URL.Path)
	respondWithJSON(w, 200, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	})
}

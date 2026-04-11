package main

import (
	"log"
	"net/http"
	"slices"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	arr := []Chirp{}
	author := r.URL.Query().Get("author_id")
	if author == "" {
		chirps, err := cfg.db.GetAllChirps(r.Context())
		if err != nil {
			respondWithError(w, 500, err.Error())
			return
		}
		for _, chirp := range chirps {

			arr = append(arr, Chirp{
				ID:        chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body:      chirp.Body,
				UserId:    chirp.UserID,
			})
		}

	} else {
		userID, err := uuid.Parse(author)
		if err != nil {
			respondWithError(w, 500, err.Error())
			return
		}
		chirps, err := cfg.db.GetChirpsFromUser(r.Context(), userID)
		if err != nil {
			respondWithError(w, 500, err.Error())
			return
		}
		for _, chirp := range chirps {

			arr = append(arr, Chirp{
				ID:        chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body:      chirp.Body,
				UserId:    chirp.UserID,
			})
		}
	}
	sortQuery := r.URL.Query().Get("sort")
	if sortQuery == "asc" {
		slices.SortFunc(arr, func(a, b Chirp) int {
			return a.CreatedAt.Compare(b.CreatedAt)
		})
	} else if sortQuery == "desc" {
		slices.SortFunc(arr, func(a, b Chirp) int {
			return b.CreatedAt.Compare(a.CreatedAt)
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

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type RequestBody struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	reqBody := RequestBody{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		respondWithError(w, 500, "Could not decode")
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), reqBody.Email)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	userMain := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	log.Printf("%s %s", r.Method, r.URL.Path)
	respondWithJSON(w, 201, userMain)

}

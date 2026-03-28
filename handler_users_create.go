package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/cfranklin121/chirpy/internal/auth"
	"github.com/cfranklin121/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Token     string    `json:"token"`
}

func (cfg *apiConfig) handlerUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type RequestBody struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type ReturnVal struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	reqBody := RequestBody{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		respondWithError(w, 500, "Could not decode")
		return
	}
	log.Println(reqBody)

	hash, err := auth.HashPassword(reqBody.Password)
	if err != nil {
		respondWithError(w, 500, "Could not hash password")
		return
	}

	params := database.CreateUserParams{
		Email:          reqBody.Email,
		HashedPassword: hash,
	}

	user, err := cfg.db.CreateUser(r.Context(), params)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	log.Printf("%s %s", r.Method, r.URL.Path)
	respondWithJSON(w, 201, ReturnVal{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})

}

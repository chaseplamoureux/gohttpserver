package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/chaseplamoureux/gohttpserver/internal/auth"
	"github.com/chaseplamoureux/gohttpserver/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	fmt.Println(params.Email)
	if err != nil {
		errorResponse(w, 500, "Couldnt decode parameters", err)
		return
	}
	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		errorResponse(w, 500, "Error occurred while hashing password", err)
	}
	newUser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: sql.NullString{String: hashedPass, Valid: true},
	})
	if err != nil {
		errorResponse(w, 500, "Could not create user in DB", err)
	}

	responseJSON(w, 201, response{
		User: User{
			ID:        newUser.ID,
			CreatedAt: newUser.CreatedAt,
			UpdatedAt: newUser.UpdatedAt,
			Email:     newUser.Email,
		},
	})
}

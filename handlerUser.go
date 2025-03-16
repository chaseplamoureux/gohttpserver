package main

import (
	"context"
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
	ID         uuid.UUID `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Email      string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
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
		return
	}
	newUser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: sql.NullString{String: hashedPass, Valid: true},
	})
	if err != nil {
		errorResponse(w, 500, "Could not create user in DB", err)
		return
	}

	responseJSON(w, 201, response{
		User: User{
			ID:         newUser.ID,
			CreatedAt:  newUser.CreatedAt,
			UpdatedAt:  newUser.UpdatedAt,
			Email:      newUser.Email,
			IsChirpyRed: newUser.IsChirpyRed.Bool,
		},
	})
}

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		errorResponse(w, 401, "No bearer auth token in request headers", err)
	}

	userId, err := auth.ValidateJWT(bearerToken, cfg.jwt)
	if err != nil {
		errorResponse(w, 401, "Could not validate JWT", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		errorResponse(w, 500, "Couldnt decode parameters", err)
		return
	}
	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		errorResponse(w, 500, "Error occurred while hashing password", err)
		return
	}

	updatedUser, err := cfg.db.UpdateUserEmailandPassword(context.Background(), database.UpdateUserEmailandPasswordParams{
		Email:          params.Email,
		HashedPassword: sql.NullString{String: hashedPass, Valid: true},
		ID:             userId,
	})
	if err != nil {
		errorResponse(w, 500, "Could not create user in DB", err)
		return
	}
	responseJSON(w, 200, response{
		User: User{
			ID:         updatedUser.ID,
			CreatedAt:  updatedUser.CreatedAt,
			UpdatedAt:  updatedUser.UpdatedAt,
			Email:      updatedUser.Email,
			IsChirpyRed: updatedUser.IsChirpyRed.Bool,
		},
	})

}

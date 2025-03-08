package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/chaseplamoureux/gohttpserver/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		errorResponse(w, 500, "could not decode parameters", err)
	}
	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		errorResponse(w, 500, "Could not get user details", err)
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword.String)
	if err != nil {
		errorResponse(w, 401, "Incorrect email or password", err)
	}

	expirationTime := time.Hour
	if params.ExpiresInSeconds > 0 && params.ExpiresInSeconds < 3600 {
		expirationTime = time.Duration(params.ExpiresInSeconds) * time.Second
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.jwt, expirationTime)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Couldn't create access JWT", err)
		return
	}
	responseJSON(w, 200, response{User: User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	},
		Token: accessToken,
	})

}

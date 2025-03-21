package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/chaseplamoureux/gohttpserver/internal/auth"
	"github.com/chaseplamoureux/gohttpserver/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerUserChirp(w http.ResponseWriter, r *http.Request) {
	// this function needs to take in a post request and process the data in the body.
	type parameters struct {
		Body string `json:"body"`
	}

	type response struct {
		Chirp
	}
	jwt, err := auth.GetBearerToken(r.Header)

	if err != nil {
		errorResponse(w, 401, "No auth token present", err)
		return
	}
	fmt.Printf("JWT: %s\n", jwt)
	userId, err := auth.ValidateJWT(jwt, cfg.jwt)
	if err != nil {
		errorResponse(w, 401, "Could not validate JWT", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	fmt.Println(params.Body)
	if err != nil {
		errorResponse(w, 500, "Couldnt decode parameters", err)
		return
	}

	if len(params.Body) > 0 && len(params.Body) <= 140 {
		params.Body = profanityCheck(params.Body)
		dbParams := database.CreateChirpParams{Body: params.Body, UserID: userId}
		dbResponse, err := cfg.db.CreateChirp(r.Context(), dbParams)
		if err != nil {
			errorResponse(w, 500, "Could not save new chirp", err)
		}
		responseJSON(w, 201, response{Chirp: Chirp{
			Id:        dbResponse.ID,
			CreatedAt: dbResponse.CreatedAt,
			UpdatedAt: dbResponse.UpdatedAt,
			Body:      dbResponse.Body,
			UserId:    dbResponse.UserID,
		}})
		return
	} else if len(params.Body) > 140 {
		errorResponse(w, 400, "Chirp is too long. Must be less than 140 characters", nil)
	} else {
		errorResponse(w, 500, "Error has occurred", nil)
	}

}

func profanityCheck(body string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(body, " ")

	for i, word := range words {
		if slices.Contains(badWords, strings.ToLower(word)) {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		errorResponse(w, 500, "Could not retreive chirps from server", err)
		return
	}
	response := []Chirp{}

	for _, dbChirp := range dbChirps {
		response = append(response, Chirp{
			Id:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserId:    dbChirp.UserID,
		})
	}
	responseJSON(w, 200, response)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpId, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Provided ID is not a valid UUID", err)
	}
	fmt.Println(chirpId)
	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpId)
	if err != nil {
		errorResponse(w, 404, "chirp with ID not found", err)
		return
	}
	responseJSON(w, 200, Chirp{
		Id:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserId:    dbChirp.UserID,
	})
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		errorResponse(w, 401, "No bearer auth token in request headers", err)
		return
	}

	userId, err := auth.ValidateJWT(bearerToken, cfg.jwt)
	if err != nil {
		errorResponse(w, 401, "Could not validate JWT", err)
		return
	}
	chirpId, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Provided ID is not a valid UUID", err)
		return
	}
	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpId)
	if err != nil {
		errorResponse(w, 404, "chirp with ID not found", err)
		return
	}

	if dbChirp.UserID != userId {
		errorResponse(w, 403, "Invalid permissions", err)
		return
	}
	err = cfg.db.DeleteChirp(context.Background(), dbChirp.ID)
	if err != nil {
		errorResponse(w, 404, "chirp not found", err)
		return
	}
	responseJSON(w, 204, nil)

}

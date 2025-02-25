package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

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
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}

	type response struct {
		Chirp
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	fmt.Println(params.Body)
	if err != nil {
		errorResponse(w, 500, "Couldnt decode parameters", err)
		return
	}

	if len(params.Body) > 0 && len(params.Body) <= 140 {
		params.Body = profanityCheck(params.Body)
		dbParams := database.CreateChirpParams{Body: params.Body, UserID: params.UserId}
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

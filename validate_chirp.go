package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
)

func (cfg *apiConfig) handlerChirp(w http.ResponseWriter, r *http.Request) {
	// this function needs to take in a post request and process the data in the body.
	type chirp struct {
		Body string `json:"body"`
	}

	type success struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := chirp{}
	err := decoder.Decode(&params)
	fmt.Println(params.Body)
	if err != nil {
		errorResponse(w, 500, "Couldnt decode paramters", err)
		return
	}

	if len(params.Body) > 0 && len(params.Body) <= 140 {
		params.Body = profanityCheck(params.Body)
		responseJSON(w, 200, success{CleanedBody: params.Body})
		return
	}
	errorResponse(w, 400, "Chirp is too long", nil)

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
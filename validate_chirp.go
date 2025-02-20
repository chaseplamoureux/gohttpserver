package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerChirp(w http.ResponseWriter, r *http.Request) {
	// this function needs to take in a post request and process the data in the body.
	type chirp struct {
		Body string `json:"body"`
	}

	type success struct {
		Valid bool `json:"valid"`
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
		responseJSON(w, 200, success{Valid: true})
		return
	}
	errorResponse(w, 400, "Chirp is too long", nil)

}
package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func errorResponse(w http.ResponseWriter, statusCode int, errorMsg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if statusCode > 499 {
		log.Printf("Respondng with 5XX error %s", errorMsg)
	}
	type httpError struct {
		Error string `json:"error"`
	}

	failureReason := httpError{
		Error: errorMsg,
	}

	responseJSON(w, statusCode, failureReason)

}

func responseJSON(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshalling json %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(dat)
}

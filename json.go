package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {

	data, err := json.Marshal(payload)

	if err != nil {
		log.Printf("Something went wrong during parsing json, %v", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)

}

func respondWithError(w http.ResponseWriter, code int, msg string) {

	if code > 499 {
		log.Printf("Response with Error: %s", msg)
	}
	type errResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errResponse{Error: msg})
}

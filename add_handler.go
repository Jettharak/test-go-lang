package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

type addRequest struct {
	Num1 int `json:"num1"`
	Num2 int `json:"num2"`
}

type addResponse struct {
	Result int `json:"result"`
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request addRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		http.Error(w, "request body must contain valid JSON with integer fields num1 and num2", http.StatusBadRequest)
		return
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		http.Error(w, "request body must contain exactly one JSON object", http.StatusBadRequest)
		return
	}

	log.Printf("Received request: %+v", request)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(addResponse{Result: add(request.Num1, request.Num2)}); err != nil {
		log.Printf("encoding response: %v", err)
	}
}

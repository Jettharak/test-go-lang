package main

import (
	"log"
	"net/http"
)

func startServer() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/add", addHandler)
	log.Println("API server listening on :8080")
	return http.ListenAndServe(":8080", mux)
}

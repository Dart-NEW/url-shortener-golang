package main

import (
	"log"
	"net/http"

	"url-shortener-golang/api"
	"url-shortener-golang/storage"
)

func main() {
	st := storage.NewMemoryStorage()

	h := api.NewHandler(st)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

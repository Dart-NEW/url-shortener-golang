package main

import (
	"flag"
	"log"
	"net/http"

	"url-shortener-golang/api"
	"url-shortener-golang/storage"
)

func main() {
	storageType := flag.String("storage", "memory", "Storage type (memory|postgres)")
	postgresConn := flag.String("postgres-conn", "", "PostgreSQL connection string")
	flag.Parse()
	var st storage.Storage
	var err error
	switch *storageType {
	case "memory":
		st = storage.NewMemoryStorage()
		log.Println("Using in-memory storage")
	case "postgres":
		if *postgresConn == "" {
			log.Fatal("PostgreSQL connection string is required")
		}
		st, err = storage.NewPostgresStorage(*postgresConn)
		if err != nil {
			log.Fatalf("Failed to initialize PostgreSQL storage: %v", err)
		}
		defer st.Close()
		log.Println("Using PostgreSQL storage")
	default:
		log.Fatalf("unsupported storage type: %s", *storageType)
	}

	h := api.NewHandler(st)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

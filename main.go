package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)	
	}

	router := chi.NewRouter()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome"))
	})

	listenAddr := os.Getenv("LISTEN_ADDR")
	log.Println("Listening on port", listenAddr)
	http.ListenAndServe(listenAddr, router)
}

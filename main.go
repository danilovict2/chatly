package main

import (
	"log"
	"net/http"
	"os"

	"github.com/danilovict2/go-real-time-chat/controllers"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	router := chi.NewRouter()
	router.Get("/", controllers.Make(controllers.HomeIndex))

	listenAddr := os.Getenv("LISTEN_ADDR")
	log.Println("Listening on port", listenAddr)
	http.ListenAndServe(listenAddr, router)
}

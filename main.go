package main

import (
	"log"
	"net/http"
	"os"

	"github.com/danilovict2/go-real-time-chat/controllers"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type Server struct {
	router chi.Router
}

func NewServer() Server {
	return Server{
		router: router(),
	}
}

func router() chi.Router {
	r := chi.NewRouter()
	r.Handle("/public/*", http.StripPrefix("/public/", http.FileServerFS(os.DirFS("/public/"))))

	// Auth routes
	r.Group(func(r chi.Router) {
		r.Get("/register", controllers.Make(controllers.RegisterForm))
		r.Post("/register", controllers.Make(controllers.Register))
	})

	r.Get("/", controllers.Make(controllers.HomeIndex))

	return r
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	server := NewServer()
	server.ListenAndServe()
}

func (s *Server) ListenAndServe() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	log.Println("Listening on port", listenAddr)
	http.ListenAndServe(listenAddr, s.router)
}

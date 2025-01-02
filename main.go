package main

import (
	"log"
	"net/http"
	"os"

	"github.com/danilovict2/go-real-time-chat/controllers"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	router   chi.Router
	database *gorm.DB
}

func NewServer() (Server, error) {
	db, err := database()
	if err != nil {
		return Server{}, err
	}

	r := router()

	return Server{
		router:   r,
		database: db,
	}, nil
}

func router() chi.Router {
	r := chi.NewRouter()
	r.Handle("/public/*", http.StripPrefix("/public/", http.FileServerFS(os.DirFS("/public/"))))

	// Auth routes
	r.Group(func(r chi.Router) {
		r.Get("/register", controllers.Make(controllers.RegisterForm))
	})

	r.Get("/", controllers.Make(controllers.HomeIndex))

	return r
}

func database() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	return db, err
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	server, err := NewServer()
	if err != nil {
		log.Fatal(err)
	}
	
	server.ListenAndServe()
}

func (s *Server) ListenAndServe() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	log.Println("Listening on port", listenAddr)
	http.ListenAndServe(listenAddr, s.router)
}
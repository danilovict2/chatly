package main

import (
	"log"
	"net/http"
	"os"

	"github.com/danilovict2/go-real-time-chat/controllers"
	"github.com/danilovict2/go-real-time-chat/jwt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
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
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(controllers.UserFromJWTMiddleware)

	r.Handle("/public/*", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	tokenAuth := jwt.NewAuth()

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(controllers.Authenticator("/login"))

		r.Get("/protected", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Protected route"))
		})
	})

	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/", controllers.Make(controllers.HomeIndex))

		// Auth routes
		r.Group(func(r chi.Router) {
			r.Get("/register", controllers.Make(controllers.RegisterForm))
			r.Post("/register", controllers.Make(controllers.Register))

			r.Get("/login", controllers.Make(controllers.LoginForm))
			r.Post("/login", controllers.Make(controllers.Login))

			r.Get("/logout", controllers.Make(controllers.Logout))
		})
	})

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

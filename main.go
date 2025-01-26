package main

import (
	"crypto/rand"
	"log"
	"net/http"
	"os"

	"github.com/danilovict2/go-real-time-chat/controllers"
	"github.com/danilovict2/go-real-time-chat/internal/jwt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/gorilla/csrf"
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
	r.Use(csrf.Protect(mustGenerateCSRFKey(), csrf.Path("/")))

	r.Handle("/public/*", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	tokenAuth := jwt.NewAuth()

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(controllers.Authenticator("/login"))

		r.Get("/profile", controllers.Make(controllers.ProfileShow))
		r.Post("/profile/update", controllers.Make(controllers.ProfileUpdate))

		r.Get("/", controllers.Make(controllers.ChatShow))
		r.Get("/chat/{receiverUsername}", controllers.Make(controllers.ChatShow))

		r.Post("/message/{receiverUsername}", controllers.Make(controllers.MessageStore))
	})

	// Public routes
	r.Group(func(r chi.Router) {
		// Auth routes
		r.Group(func(r chi.Router) {
			r.Get("/register", controllers.Make(controllers.RegisterForm))
			r.Post("/register", controllers.Make(controllers.Register))

			r.Get("/login", controllers.Make(controllers.LoginForm))
			r.Post("/login", controllers.Make(controllers.Login))

			r.Post("/logout", controllers.Make(controllers.Logout))
		})
	})

	return r
}

func mustGenerateCSRFKey() (key []byte) {
	key = make([]byte, 32)
	n, err := rand.Read(key)
	if err != nil {
		panic(err)
	}
	if n != 32 {
		panic("unable to read 32 bytes for CSRF key")
	}
	return
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

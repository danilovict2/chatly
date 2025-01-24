package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/danilovict2/go-real-time-chat/internal/database"
	"github.com/danilovict2/go-real-time-chat/internal/jwt"
	"github.com/danilovict2/go-real-time-chat/models"
	"github.com/go-chi/jwtauth/v5"
)

type ControllerError struct {
	err  error
	code int
}

type HTTPController func(w http.ResponseWriter, r *http.Request) ControllerError

func Make(h HTTPController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if controllerErr := h(w, r); controllerErr != (ControllerError{}) {
			log.Println("HTTP controller error:", controllerErr.err, "Path:", r.URL.Path, "Status code:", controllerErr.code)

			w.WriteHeader(controllerErr.code)
			if controllerErr.code >= 500 {
				_, err := w.Write([]byte("Whoops. Something went wrong."))
				if err != nil {
					log.Println("Write error:", err)
					return
				}
			} else if controllerErr.code >= 400 {
				out := fmt.Sprintf("There was an error with your request: %v", controllerErr.err)
				_, err := w.Write([]byte(out))
				if err != nil {
					log.Println("Write error:", err)
					return
				}
			}
		}
	}
}

func Render(w http.ResponseWriter, r *http.Request, c templ.Component) ControllerError {
	return ControllerError{
		err:  c.Render(r.Context(), w),
		code: http.StatusInternalServerError,
	}
}

func Authenticator(loginRoute string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			token, _, err := jwtauth.FromContext(r.Context())

			if err != nil {
				http.Redirect(w, r, loginRoute, http.StatusUnauthorized)
				return
			}

			if token == nil {
				http.Redirect(w, r, loginRoute, http.StatusUnauthorized)
				return
			}

			// Token is authenticated, pass it through
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(hfn)
	}
}

type contextKey string

const userContextKey contextKey = "user"

func UserFromJWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := jwtauth.TokenFromCookie(r)
		ja := jwt.NewAuth()
		token, err := ja.Decode(tokenString)
		if err != nil || token == nil {
			next.ServeHTTP(w, r)
			return
		}

		userID, ok := token.Get("user_id")
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		db, err := database.NewConnection()
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		user := &models.User{}
		if err := db.First(user, userID).Error; err != nil {
			fmt.Println(err)
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

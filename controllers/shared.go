package controllers

import (
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/jwtauth/v5"
)

type HTTPController func(w http.ResponseWriter, r *http.Request) error

func Make(h HTTPController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			log.Fatal("HTTP controller error", err, "Path", r.URL.Path)
		}
	}
}

func Render(w http.ResponseWriter, r *http.Request, c templ.Component) error {
	return c.Render(r.Context(), w)
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
package controllers

import (
	"log"
	"net/http"

	"github.com/a-h/templ"
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

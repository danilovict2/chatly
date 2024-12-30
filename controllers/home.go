package controllers

import (
	"github.com/danilovict2/go-real-time-chat/views/home"
	"net/http"
)

func HomeIndex(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, home.Home())
}
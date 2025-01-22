package controllers

import (
	"net/http"

	"github.com/danilovict2/go-real-time-chat/models"
	"github.com/danilovict2/go-real-time-chat/views/home"
)

func HomeIndex(w http.ResponseWriter, r *http.Request) error {
	user, _ := r.Context().Value(userContextKey).(*models.User)
	
	return Render(w, r, home.Home(user))
}
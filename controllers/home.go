package controllers

import (
	"net/http"

	"github.com/danilovict2/go-real-time-chat/internal/database"
	"github.com/danilovict2/go-real-time-chat/models"
	"github.com/danilovict2/go-real-time-chat/views/home"
)

func HomeIndex(w http.ResponseWriter, r *http.Request) ControllerError {
	authUser, _ := r.Context().Value(userContextKey).(*models.User)
	db, err := database.NewConnection()
	if err != nil {
		return ControllerError{
			err:  err,
			code: http.StatusInternalServerError,
		}
	}

	users := make([]models.User, 0)
	if err := db.Find(&users).Error; err != nil {
		return ControllerError{
			err:  err,
			code: http.StatusInternalServerError,
		}
	}

	return Render(w, r, home.Home(authUser, users))
}

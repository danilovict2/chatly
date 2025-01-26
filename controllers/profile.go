package controllers

import (
	"net/http"

	"github.com/danilovict2/go-real-time-chat/internal/database"
	"github.com/danilovict2/go-real-time-chat/models"
	"github.com/danilovict2/go-real-time-chat/views/profile"
)

func ProfileShow(w http.ResponseWriter, r *http.Request) ControllerError {
	user, _ := r.Context().Value(userContextKey).(*models.User)
	return Render(w, r, profile.Profile(*user))
}

func ProfileUpdate(w http.ResponseWriter, r *http.Request) ControllerError {
	avatarPath, controllerErr := SaveFormFile(r, "avatar")
	if controllerErr != (ControllerError{}) {
		return controllerErr
	}

	user, _ := r.Context().Value(userContextKey).(*models.User)
	user.Avatar = &avatarPath

	db, err := database.NewConnection()
	if err != nil {
		return ControllerError{
			err:  err,
			code: http.StatusInternalServerError,
		}
	}

	if err := db.Save(user).Error; err != nil {
		return ControllerError{
			err:  err,
			code: http.StatusInternalServerError,
		}
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
	return ControllerError{}
}

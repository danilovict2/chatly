package controllers

import (
	"net/http"
	"os"

	"github.com/danilovict2/go-real-time-chat/models"
	"github.com/danilovict2/go-real-time-chat/views/profile"
	"github.com/go-chi/chi/v5"
)

func (cfg *Config) ProfileShow(w http.ResponseWriter, r *http.Request) ControllerError {
	user, _ := r.Context().Value(userContextKey).(*models.User)
	username := chi.URLParam(r, "username")
	if user.Username != username {
		w.WriteHeader(http.StatusUnauthorized)
		return ControllerError{}
	}

	return Render(w, r, profile.Profile(*user))
}

func (cfg *Config) ProfileUpdate(w http.ResponseWriter, r *http.Request) ControllerError {
	user, _ := r.Context().Value(userContextKey).(*models.User)
	username := chi.URLParam(r, "username")
	if user.Username != username {
		w.WriteHeader(http.StatusUnauthorized)
		return ControllerError{}
	}

	avatarPath, controllerErr := SaveFormFile(r, "avatar")
	if controllerErr != (ControllerError{}) {
		return controllerErr
	}

	if user.Avatar != nil {
		err := os.Remove("." + os.Getenv("IMG_ROOT") + *user.Avatar)
		if err != nil {
			return ControllerError{
				err:  err,
				code: http.StatusInternalServerError,
			}
		}
	}

	user.Avatar = &avatarPath

	if err := cfg.DB.Save(user).Error; err != nil {
		return ControllerError{
			err:  err,
			code: http.StatusInternalServerError,
		}
	}

	http.Redirect(w, r, "/profile/"+username, http.StatusSeeOther)
	return ControllerError{}
}

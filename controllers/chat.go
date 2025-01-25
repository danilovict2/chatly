package controllers

import (
	"net/http"

	"github.com/danilovict2/go-real-time-chat/internal/database"
	"github.com/danilovict2/go-real-time-chat/models"
	"github.com/danilovict2/go-real-time-chat/views/chat"
	"github.com/go-chi/chi/v5"
)

func ChatShow(w http.ResponseWriter, r *http.Request) ControllerError {
	authUser, _ := r.Context().Value(userContextKey).(*models.User)
	var selectedUser *models.User = nil
	
	if selectedUserUsername := chi.URLParam(r, "username"); selectedUserUsername != "" {
		db, err := database.NewConnection()
		if err != nil {
			return ControllerError{
				err: err,
				code: http.StatusInternalServerError,
			}
		}
		
		selectedUser = &models.User{}
		if err := db.First(selectedUser, selectedUserUsername).Error; err != nil {
			return ControllerError{
				err: err,
				code: http.StatusInternalServerError,
			}
		}
	}

	return Render(w, r, chat.Chat(authUser, selectedUser))
}

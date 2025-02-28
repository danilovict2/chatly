package controllers

import (
	"github.com/danilovict2/go-real-time-chat/internal/repository"
	"net/http"

	"github.com/danilovict2/go-real-time-chat/models"
	"github.com/danilovict2/go-real-time-chat/views/chat"
	"github.com/go-chi/chi/v5"
)

func (cfg *Config) ChatShow(w http.ResponseWriter, r *http.Request) ControllerError {
	sender, _ := r.Context().Value(userContextKey).(*models.User)
	var receiver *models.User = nil
	messages := make([]models.Message, 0)

	if receiverUsername := chi.URLParam(r, "receiverUsername"); receiverUsername != "" {
		receiver = &models.User{}
		if err := cfg.DB.Where("username = ?", receiverUsername).First(receiver).Error; err != nil {
			return ControllerError{
				err:  err,
				code: http.StatusInternalServerError,
			}
		}

		ids := []uint{sender.ID, receiver.ID}
		if err := cfg.DB.Where("sender_id IN ? AND receiver_id IN ?", ids, ids).Find(&messages).Error; err != nil {
			return ControllerError{
				err:  err,
				code: http.StatusInternalServerError,
			}
		}
	}

	ur := repository.NewUserRepository(cfg.DB)
	users, err := ur.FindAllExcept(*sender)
	if err != nil {
		return ControllerError{
			err:  err,
			code: http.StatusInternalServerError,
		}
	}

	return Render(w, r, chat.Chat(sender, receiver, messages, users))
}

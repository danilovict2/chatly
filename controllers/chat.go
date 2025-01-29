package controllers

import (
	"net/http"

	"github.com/danilovict2/go-real-time-chat/internal/database"
	"github.com/danilovict2/go-real-time-chat/internal/pusher"
	"github.com/danilovict2/go-real-time-chat/models"
	"github.com/danilovict2/go-real-time-chat/views/chat"
	"github.com/danilovict2/go-real-time-chat/views/components"
	"github.com/go-chi/chi/v5"
)

func ChatShow(w http.ResponseWriter, r *http.Request) ControllerError {
	sender, _ := r.Context().Value(userContextKey).(*models.User)
	var receiver *models.User = nil
	messages := make([]models.Message, 0)

	if receiverUsername := chi.URLParam(r, "receiverUsername"); receiverUsername != "" {
		db, err := database.NewConnection()
		if err != nil {
			return ControllerError{
				err:  err,
				code: http.StatusInternalServerError,
			}
		}

		receiver = &models.User{}
		if err := db.Where("username = ?", receiverUsername).First(receiver).Error; err != nil {
			return ControllerError{
				err:  err,
				code: http.StatusInternalServerError,
			}
		}

		ids := []uint{sender.ID, receiver.ID}
		if err := db.Where("sender_id IN ? AND receiver_id IN ?", ids, ids).Find(&messages).Error; err != nil {
			return ControllerError{
				err:  err,
				code: http.StatusInternalServerError,
			}
		}
	}

	return Render(w, r, chat.Chat(sender, receiver, messages))
}

func MessageStore(w http.ResponseWriter, r *http.Request) ControllerError {
	messageImage, controllerErr := SaveFormFile(r, "message_image")
	if controllerErr != (ControllerError{}) {
		return controllerErr
	}

	receiverUsername := chi.URLParam(r, "receiverUsername")
	db, err := database.NewConnection()
	if err != nil {
		return ControllerError{
			err:  err,
			code: http.StatusInternalServerError,
		}
	}

	receiver := &models.User{}
	if err := db.Where("username = ?", receiverUsername).First(receiver).Error; err != nil {
		return ControllerError{
			err:  err,
			code: http.StatusInternalServerError,
		}
	}

	sender, _ := r.Context().Value(userContextKey).(*models.User)

	message := models.Message{
		SenderID:   sender.ID,
		ReceiverID: receiver.ID,
		Text:       r.PostFormValue("message"),
		Image:      messageImage,
	}

	if err := db.Create(&message).Error; err != nil {
		return ControllerError{
			err:  err,
			code: http.StatusInternalServerError,
		}
	}

	data := map[string]string{
		"Text":         message.Text,
		"Image":        message.Image,
		"Sender":       sender.Username,
		"SenderAvatar": components.Avatar(*sender),
		"CreatedAt":    message.CreatedAt.Format("3:04 PM"),
	}

	pusherClient := pusher.NewClient()
	pusherClient.Trigger("message", "to."+receiverUsername, data)

	http.Redirect(w, r, "/chat/"+receiver.Username, http.StatusSeeOther)
	return ControllerError{}
}

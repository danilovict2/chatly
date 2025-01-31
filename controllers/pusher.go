package controllers

import (
	"io"
	"net/http"
	"strconv"

	pusherClient "github.com/danilovict2/go-real-time-chat/internal/pusher"
	"github.com/danilovict2/go-real-time-chat/models"
	"github.com/pusher/pusher-http-go/v5"
)

func PusherAuth(w http.ResponseWriter, r *http.Request) ControllerError {
	params, err := io.ReadAll(r.Body)
	if err != nil {
		return ControllerError{
			err:  err,
			code: http.StatusForbidden,
		}
	}

	user, _ := r.Context().Value(userContextKey).(*models.User)
	member := pusher.MemberData{
		UserID: strconv.Itoa(int(user.ID)),
	}

	pusherClient := pusherClient.NewClient()
	response, err := pusherClient.AuthorizePresenceChannel(params, member)
	if err != nil {
		return ControllerError{
			err: err,
			code: http.StatusUnauthorized,
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
	
	return ControllerError{}
}

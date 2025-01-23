package controllers

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/danilovict2/go-real-time-chat/internal/database"
	"github.com/danilovict2/go-real-time-chat/models"
	"github.com/danilovict2/go-real-time-chat/views/profile"
)

func ProfileShow(w http.ResponseWriter, r *http.Request) error {
	user, _ := r.Context().Value(userContextKey).(*models.User)
	return Render(w, r, profile.Profile(*user))
}

func ProfileUpdate(w http.ResponseWriter, r *http.Request) error {
	const maxMemory = 1 << 20
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		return err
	}
	
	user, _ := r.Context().Value(userContextKey).(*models.User)
	avatar, header, err := r.FormFile("avatar")
	if err != nil {
		return err
	}
	defer avatar.Close()
	
	allowedMimeTypes := []string{"image/png", "image/jpg"}
	mimeType, _, err := mime.ParseMediaType(header.Header.Get("Content-Type"))
	if err != nil || !slices.Contains(allowedMimeTypes, mimeType) {
		return fmt.Errorf("invalid file extension")	
	}

	ext := strings.Split(mimeType, "/")[1]
	fName := fmt.Sprintf("%s%s.%s", os.Getenv("IMG_ROOT"), user.Username, ext)

	file, err := os.Create(fName)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.Copy(file, avatar); err != nil {
		return err
	}

	user.Avatar = &fName

	db, err := database.NewConnection()
	if err != nil {
		return err
	}

	if err := db.Save(user).Error; err != nil {
		return err
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
	return nil
}
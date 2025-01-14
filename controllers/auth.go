package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/danilovict2/go-real-time-chat/internal/database"
	"github.com/danilovict2/go-real-time-chat/jwt"
	"github.com/danilovict2/go-real-time-chat/models"
	"github.com/danilovict2/go-real-time-chat/views/auth"
	"golang.org/x/crypto/bcrypt"
)

func RegisterForm(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, auth.Register())
}

func Register(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	db, err := database.NewConnection()
	if err != nil {
		return err
	}

	password, err := bcrypt.GenerateFromPassword([]byte(r.PostFormValue("password")), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{
		Username: r.PostFormValue("username"),
		Email:    r.PostFormValue("email"),
		Password: password,
	}

	if valid, err := user.IsValid(db); !valid {
		// TODO: Implement error handling
		fmt.Println(err)
		http.Redirect(w, r, "/", http.StatusFound)
		return nil
	}

	db.Create(&user)

	tokenAuth := jwt.NewAuth()
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"user_id": user.ID})
	if err != nil {
		return err
	}

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		Secure:   true,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/protected", http.StatusFound)

	return nil
}

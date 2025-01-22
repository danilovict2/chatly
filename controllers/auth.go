package controllers

import (
	"net/http"
	"time"

	"github.com/danilovict2/go-real-time-chat/internal/database"
	"github.com/danilovict2/go-real-time-chat/jwt"
	"github.com/danilovict2/go-real-time-chat/models"
	"github.com/danilovict2/go-real-time-chat/views/auth"
	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"
)

const DefaultJWTExpiration time.Duration = time.Hour * 24 * 7

func RegisterForm(w http.ResponseWriter, r *http.Request) error {
	token := jwtauth.TokenFromCookie(r)
	// Prevent authenticated users from accesing register
	if token != "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return nil
	}

	errorMessage := r.URL.Query().Get("error_message")
	return Render(w, r, auth.Register(errorMessage))
}

func Register(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	db, err := database.NewConnection()
	if err != nil {
		return err
	}

	user := models.User{
		Username: r.PostFormValue("username"),
		Email:    r.PostFormValue("email"),
		Password: []byte(r.PostFormValue("password")),
	}

	if valid, reason := user.IsValid(db); !valid {
		http.Redirect(w, r, "/register?error_message="+reason, http.StatusFound)
		return nil
	}

	user.Password, err = bcrypt.GenerateFromPassword(user.Password, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	db.Create(&user)
	err = setJWTCookie(user.ID, time.Now().Add(DefaultJWTExpiration), w)
	if err != nil {
		return err
	}

	http.Redirect(w, r, "/protected", http.StatusFound)
	return nil
}

func LoginForm(w http.ResponseWriter, r *http.Request) error {
	token := jwtauth.TokenFromCookie(r)
	// Prevent authenticated users from accesing login
	if token != "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return nil
	}

	return Render(w, r, auth.Login())
}

func Login(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	db, err := database.NewConnection()
	if err != nil {
		return err
	}

	user := models.User{}
	if err := db.Where("email = ?", r.PostFormValue("email")).First(&user).Error; err != nil {
		// TODO: Implement error handling
		return err
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(r.PostFormValue("password"))); err != nil {
		// TODO: Implement error handling
		return err
	}

	err = setJWTCookie(user.ID, time.Now().Add(DefaultJWTExpiration), w)
	if err != nil {
		return err
	}

	http.Redirect(w, r, "/protected", http.StatusFound)
	return nil
}

func setJWTCookie(userID uint, expires time.Time, w http.ResponseWriter) error {
	tokenAuth := jwt.NewAuth()
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"user_id": userID})
	if err != nil {
		return err
	}

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  expires,
		Secure:   true,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
	return nil
}

func Logout(w http.ResponseWriter, r *http.Request) error {
	setJWTCookie(0, time.Now(), w)
	http.Redirect(w, r, "/", http.StatusSeeOther)

	return nil
}

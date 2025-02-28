package controllers

import (
	"net/http"
	"time"

	"github.com/danilovict2/go-real-time-chat/internal/jwt"
	"github.com/danilovict2/go-real-time-chat/models"
	"github.com/danilovict2/go-real-time-chat/views/auth"
	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"
)

const DefaultJWTExpiration time.Duration = time.Hour * 24 * 7

func (cfg *Config) RegisterForm(w http.ResponseWriter, r *http.Request) ControllerError {
	token := jwtauth.TokenFromCookie(r)
	// Prevent authenticated users from accesing register
	if token != "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return ControllerError{}
	}

	errorMessage := r.URL.Query().Get("error_message")
	return Render(w, r, auth.Register(errorMessage))
}

func (cfg *Config) Register(w http.ResponseWriter, r *http.Request) ControllerError {
	if err := r.ParseForm(); err != nil {
		return ControllerError{
			err:  err,
			code: http.StatusBadRequest,
		}
	}

	user := models.User{
		Username: r.PostFormValue("username"),
		Email:    r.PostFormValue("email"),
		Password: []byte(r.PostFormValue("password")),
	}

	if valid, reason := user.IsValid(cfg.DB); !valid {
		http.Redirect(w, r, "/register?error_message="+reason, http.StatusFound)
		return ControllerError{}
	}

	var err error
	user.Password, err = bcrypt.GenerateFromPassword(user.Password, bcrypt.DefaultCost)
	if err != nil {
		return ControllerError{
			err:  err,
			code: http.StatusInternalServerError,
		}
	}

	cfg.DB.Create(&user)
	if err = setJWTCookie(user.ID, time.Now().Add(DefaultJWTExpiration), w); err != nil {
		return ControllerError{
			err:  err,
			code: http.StatusBadRequest,
		}
	}

	http.Redirect(w, r, "/", http.StatusFound)
	return ControllerError{}
}

func (cfg *Config) LoginForm(w http.ResponseWriter, r *http.Request) ControllerError {
	token := jwtauth.TokenFromCookie(r)
	// Prevent authenticated users from accesing login
	if token != "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return ControllerError{}
	}

	return Render(w, r, auth.Login())
}

func (cfg *Config) Login(w http.ResponseWriter, r *http.Request) ControllerError {
	if err := r.ParseForm(); err != nil {
		return ControllerError{
			err:  err,
			code: http.StatusBadRequest,
		}
	}

	user := models.User{}
	if err := cfg.DB.Where("email = ?", r.PostFormValue("email")).First(&user).Error; err != nil {
		return ControllerError{
			err:  err,
			code: http.StatusInternalServerError,
		}
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(r.PostFormValue("password"))); err != nil {
		return ControllerError{
			err:  err,
			code: http.StatusBadRequest,
		}
	}

	if err := setJWTCookie(user.ID, time.Now().Add(DefaultJWTExpiration), w); err != nil {
		return ControllerError{
			err:  err,
			code: http.StatusBadRequest,
		}
	}

	http.Redirect(w, r, "/", http.StatusFound)
	return ControllerError{}
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

func (cfg *Config) Logout(w http.ResponseWriter, r *http.Request) ControllerError {
	setJWTCookie(0, time.Now(), w)
	http.Redirect(w, r, "/", http.StatusSeeOther)

	return ControllerError{}
}

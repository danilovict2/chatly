package jwt

import (
	"os"

	"github.com/go-chi/jwtauth/v5"
)

func NewAuth() *jwtauth.JWTAuth {
	jwtSecret := os.Getenv("JWT_SECRET")
	tokenAuth := jwtauth.New("HS256", []byte(jwtSecret), nil)
	
	return tokenAuth
}
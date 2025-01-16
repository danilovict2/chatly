package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex"`
	Email    string `gorm:"uniqueIndex"`
	Password []byte
}

func (user User) IsValid(db *gorm.DB) (valid bool, invalidReason string) {
	switch {
	case len(user.Username) == 0:
		return false, "Username is required"
	case len(user.Email) == 0:
		return false, "Email is required"
	case len(user.Password) < 8:
		return false, "Password must be longer than 8 characters"
	case db.Where("email = ?", user.Email).First(&User{}).Error == nil:
		return false, "User with this email already exists"
	case db.Where("username = ?", user.Username).First(&User{}).Error == nil:
		return false, "User with this username already exists"
	default:
		return true, ""
	}
}

package models

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex"`
	Email    string `gorm:"uniqueIndex"`
	Password []byte
}

func (user User) IsValid(db *gorm.DB) (valid bool, invalidReason error) {
	if err := db.Where("username = ?", user.Username).First(&User{}).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, fmt.Errorf("User with this username already exists")
	}

	if err := db.Where("email = ?", user.Email).First(&User{}).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, fmt.Errorf("User with this email already exists")
	}

	return true, nil
}
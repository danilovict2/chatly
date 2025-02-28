package repository

import (
	"github.com/danilovict2/go-real-time-chat/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) FindAllExcept(user models.User) ([]models.User, error) {
	users := make([]models.User, 0)
	if err := r.DB.Where("id <> ?", user.ID).Find(&users).Error; err != nil {

		return []models.User{}, err
	}

	return users, nil
}

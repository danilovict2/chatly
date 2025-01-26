package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	SenderID   uint
	ReceiverID uint
	Text       string `gorm:"type:text"`
	Image      string
}

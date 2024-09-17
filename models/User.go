package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserId      string `json:"user_id" gorm:"user_id"`
	DisplayName string `gorm:"uniqueIndex;not null"`
	Email       string `gorm:"uniqueIndex;not null"`
	Password    string `gorm:"not null"`
}

func (User) Tablename() string {
	return "users"
}

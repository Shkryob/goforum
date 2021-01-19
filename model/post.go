package model

import (
	"github.com/jinzhu/gorm"
)

type Post struct {
	gorm.Model

	Title    string `gorm:"not null"`
	User     User
	Body     string
	Comments []Comment
	UserID   uint
}

type Comment struct {
	gorm.Model

	Body   string `gorm:"not null"`
	PostID uint
	Name   string `gorm:"not null"`
	User   User
	UserID uint
}

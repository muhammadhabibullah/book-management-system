package models

import (
	"gorm.io/gorm"
)

// Book model
type Book struct {
	gorm.Model
	Name string `gorm:"name" json:"name" example:"The Alchemist"`
	ISBN string `gorm:"isbn" json:"isbn" example:"9780062315007"`
}

// Books model is an array of Book
type Books []Book

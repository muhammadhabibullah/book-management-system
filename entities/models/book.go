package models

import (
	"gorm.io/gorm"
)

// Book model
type Book struct {
	gorm.Model
	Name string `gorm:"name" json:"name"`
	ISBN string `gorm:"isbn" json:"isbn"`
}

// Books model is an array of Book
type Books []Book

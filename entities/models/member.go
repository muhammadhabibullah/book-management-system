package models

import (
	"gorm.io/gorm"
)

// Member model
type Member struct {
	gorm.Model
	Name string `gorm:"name" json:"name" example:"John Lennon"`
}

// Members model is an array of Member
type Members []Member

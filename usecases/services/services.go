package services

import (
	"book-management-system/repositories"
)

// Services contains services
type Services struct {
	BookService BookService
}

// Init return Services
func Init(repo *repositories.Repository) *Services {
	return &Services{
		BookService: NewBookService(repo),
	}
}

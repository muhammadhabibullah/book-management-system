package services

import (
	"book-management-system/entities/models"
	"book-management-system/repositories"
	"book-management-system/repositories/mysql"
)

// BookService handle business logic related to book
type BookService interface {
	GetBooks() (models.Books, error)
}

type bookService struct {
	BookRepository mysql.BookRepository
}

// NewBookService returns BookService
func NewBookService(repo *repositories.Repository) BookService {
	return &bookService{
		BookRepository: repo.BookRepository,
	}
}

func (svc *bookService) GetBooks() (models.Books, error) {
	return svc.BookRepository.GetAll()
}

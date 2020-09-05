package services

import (
	"book-management-system/entities/models"
	"book-management-system/repositories"
	"book-management-system/repositories/elasticsearch"
	"book-management-system/repositories/mysql"
)

// BookService handle business logic related to book
type BookService interface {
	GetBooks() (models.Books, error)
	CreateBook(*models.Book) error
	UpdateBook(*models.Book) error
}

type bookService struct {
	MySQLBookRepository mysql.BookRepository
	ESBookRepository    elasticsearch.BookRepository
}

// NewBookService returns BookService
func NewBookService(repo *repositories.Repository) BookService {
	return &bookService{
		MySQLBookRepository: repo.MySQLBookRepository,
		ESBookRepository:    repo.ESBookRepository,
	}
}

func (svc *bookService) GetBooks() (models.Books, error) {
	return svc.MySQLBookRepository.GetAll()
}

func (svc *bookService) CreateBook(book *models.Book) error {
	return svc.MySQLBookRepository.CreateBook(book)
}

func (svc *bookService) UpdateBook(book *models.Book) error {
	return svc.MySQLBookRepository.UpdateBook(book)
}

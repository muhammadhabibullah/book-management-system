package services

import (
	"log"

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
	SearchBooks(string) (models.Books, error)
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
	err := svc.MySQLBookRepository.CreateBook(book)
	if err != nil {
		return err
	}

	go func() {
		err := svc.ESBookRepository.IndexBook(book)
		if err != nil {
			log.Printf("error create book in elasticsearch %s", err)
		}
	}()
	return nil
}

func (svc *bookService) UpdateBook(book *models.Book) error {
	err := svc.MySQLBookRepository.UpdateBook(book)
	if err != nil {
		return err
	}

	go func() {
		err := svc.ESBookRepository.IndexBook(book)
		if err != nil {
			log.Printf("error update book in elasticsearch %s", err)
		}
	}()
	return nil
}

func (svc *bookService) SearchBooks(keyword string) (models.Books, error) {
	return svc.ESBookRepository.SearchBook(keyword)
}

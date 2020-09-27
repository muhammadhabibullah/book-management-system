package services

import (
	"context"
	"log"

	"book-management-system/entities/models"
	"book-management-system/repositories"
	"book-management-system/repositories/elasticsearch"
	"book-management-system/repositories/mysql"
)

// BookService handle business logic related to book
type BookService interface {
	GetBooks(context.Context) (models.Books, error)
	CreateBook(context.Context, *models.Book) error
	UpdateBook(context.Context, *models.Book) error
	SearchBooks(context.Context, string) (models.Books, error)
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

func (svc *bookService) GetBooks(ctx context.Context) (models.Books, error) {
	ctx, cancel := setContextTimeout(ctx)
	defer cancel()

	return svc.MySQLBookRepository.GetAll(ctx)
}

func (svc *bookService) CreateBook(ctx context.Context, book *models.Book) error {
	ctx, cancel := setContextTimeout(ctx)
	defer cancel()

	err := svc.MySQLBookRepository.CreateBook(ctx, book)
	if err != nil {
		return err
	}

	go func() {
		err := svc.ESBookRepository.IndexBook(context.Background(), book)
		if err != nil {
			log.Printf("error create book in elasticsearch %s", err)
		}
	}()
	return nil
}

func (svc *bookService) UpdateBook(ctx context.Context, book *models.Book) error {
	ctx, cancel := setContextTimeout(ctx)
	defer cancel()

	err := svc.MySQLBookRepository.UpdateBook(ctx, book)
	if err != nil {
		return err
	}

	go func() {
		err := svc.ESBookRepository.IndexBook(context.Background(), book)
		if err != nil {
			log.Printf("error update book in elasticsearch %s", err)
		}
	}()
	return nil
}

func (svc *bookService) SearchBooks(ctx context.Context, keyword string) (models.Books, error) {
	ctx, cancel := setContextTimeout(ctx)
	defer cancel()

	return svc.ESBookRepository.SearchBook(ctx, keyword)
}

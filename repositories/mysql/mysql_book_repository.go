package mysql

import (
	"context"

	"gorm.io/gorm"

	"book-management-system/entities/models"
)

// BookRepository handle sql query to books table
type BookRepository interface {
	GetAll(context.Context) (models.Books, error)
	CreateBook(context.Context, *models.Book) error
	UpdateBook(context.Context, *models.Book) error
}

type bookRepository struct {
	db *gorm.DB
}

// NewBookRepository returns new BookRepository
func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{
		db: db,
	}
}

func (repo *bookRepository) GetAll(ctx context.Context) (models.Books, error) {
	var books models.Books

	query := repo.db.WithContext(ctx).
		Find(&books)
	return books, query.Error
}

func (repo *bookRepository) CreateBook(ctx context.Context, book *models.Book) error {
	query := repo.db.WithContext(ctx).
		Create(book)
	return query.Error
}

func (repo *bookRepository) UpdateBook(ctx context.Context, book *models.Book) error {
	query := repo.db.WithContext(ctx).
		Updates(book)
	return query.Error
}

package mysql

import (
	"gorm.io/gorm"

	"book-management-system/entities/models"
)

// BookRepository handle sql query to book table
type BookRepository interface {
	GetAll() (models.Books, error)
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

func (repo *bookRepository) GetAll() (models.Books, error) {
	var books models.Books

	query := repo.db.Find(&books)
	return books, query.Error
}

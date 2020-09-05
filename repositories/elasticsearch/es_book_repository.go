package elasticsearch

import (
	"github.com/elastic/go-elasticsearch/v8"

	"book-management-system/entities/models"
)

// BookRepository interface
type BookRepository interface {
	CreateBook(*models.Book) error
	UpdateBook(*models.Book) error
	SearchBook(string) (models.Books, error)
}

type bookRepository struct {
	es *elasticsearch.Client
}

// NewBookRepository returns new BookRepository
func NewBookRepository(es *elasticsearch.Client) BookRepository {
	return &bookRepository{
		es: es,
	}
}

// TODO
func (repo *bookRepository) CreateBook(*models.Book) error {
	return nil
}

// TODO
func (repo *bookRepository) UpdateBook(*models.Book) error {
	return nil
}

// TODO
func (repo *bookRepository) SearchBook(string) (models.Books, error) {
	return nil, nil
}

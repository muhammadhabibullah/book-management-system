package elasticsearch

import (
	"github.com/elastic/go-elasticsearch/v8"
)

// BookRepository interface
type BookRepository interface{}

type bookRepository struct {
	es *elasticsearch.Client
}

// NewBookRepository returns new BookRepository
func NewBookRepository(es *elasticsearch.Client) BookRepository {
	return &bookRepository{
		es: es,
	}
}

package elasticsearch

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"

	"book-management-system/entities/models"
)

// BookRepository interface
type BookRepository interface {
	IndexBook(context.Context, *models.Book) error
	SearchBook(context.Context, string) (models.Books, error)
}

type bookRepository struct {
	es    *elasticsearch.Client
	index string
}

// NewBookRepository returns new BookRepository
func NewBookRepository(es *elasticsearch.Client) BookRepository {
	return &bookRepository{
		es:    es,
		index: "books",
	}
}

func (repo *bookRepository) IndexBook(ctx context.Context, book *models.Book) error {
	bookBytes, err := json.Marshal(book)
	if err != nil {
		return err
	}

	res, err := repo.es.Index(
		repo.index,
		strings.NewReader(string(bookBytes)),
		es.Index.WithContext(ctx),
		es.Index.WithDocumentID(strconv.Itoa(int(book.ID))),
		es.Index.WithPretty(),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func (repo *bookRepository) SearchBook(ctx context.Context, keyword string) (models.Books, error) {
	res, err := repo.es.Search(
		es.Search.WithContext(ctx),
		es.Search.WithIndex(repo.index),
		es.Search.WithQuery(keyword),
		es.Search.WithPretty(),
	)
	if err != nil {
		return models.Books{}, err
	}
	defer res.Body.Close()

	decodedRes := make(map[string]interface{})
	if err := json.NewDecoder(res.Body).Decode(&decodedRes); err != nil {
		return models.Books{}, err
	}

	books := make(models.Books, 0)
	for _, hit := range decodedRes["hits"].(map[string]interface{})["hits"].([]interface{}) {
		source, _ := json.Marshal(hit.(map[string]interface{})["_source"])
		var book models.Book
		err = json.Unmarshal(source, &book)
		if err != nil {
			return models.Books{}, err
		}
		books = append(books, book)
	}

	return books, nil
}

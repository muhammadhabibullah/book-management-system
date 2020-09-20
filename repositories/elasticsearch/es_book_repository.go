package elasticsearch

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"

	"book-management-system/entities/models"
)

// BookRepository interface
type BookRepository interface {
	IndexBook(*models.Book) error
	SearchBook(string) (models.Books, error)
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

func (repo *bookRepository) IndexBook(book *models.Book) error {
	bookBytes, err := json.Marshal(book)
	if err != nil {
		return err
	}

	res, err := repo.es.Index(
		repo.index,
		strings.NewReader(string(bookBytes)),
		es.Index.WithDocumentID(strconv.Itoa(int(book.ID))),
		es.Index.WithPretty(),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func (repo *bookRepository) SearchBook(keyword string) (models.Books, error) {
	res, err := repo.es.Search(
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

	var books models.Books
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

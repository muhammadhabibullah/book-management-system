// Package repositories contains all repositories
package repositories

import (
	"book-management-system/repositories/elasticsearch"
	"book-management-system/repositories/mysql"
)

// Repository contains repositories
type Repository struct {
	MySQLBookRepository mysql.BookRepository
	ESBookRepository    elasticsearch.BookRepository
}

// Init returns Repository
func Init() *Repository {
	mysqlDB := mysql.Init()
	es := elasticsearch.Init()
	return &Repository{
		MySQLBookRepository: mysql.NewBookRepository(mysqlDB),
		ESBookRepository:    elasticsearch.NewBookRepository(es),
	}
}

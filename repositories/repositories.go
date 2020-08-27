package repositories

import (
	"book-management-system/repositories/mysql"
)

// Repository contains repositories
type Repository struct {
	BookRepository mysql.BookRepository
}

// Init return Repository
func Init() *Repository {
	mysqlDB := mysql.Init()
	return &Repository{
		BookRepository: mysql.NewBookRepository(mysqlDB),
	}
}

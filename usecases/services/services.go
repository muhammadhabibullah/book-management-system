// Package services contains all services
package services

import (
	"book-management-system/repositories"
)

// Services contains services
type Services struct {
	BookService   BookService
	MemberService MemberService
}

// Init return Services
func Init(repo *repositories.Repository) *Services {
	return &Services{
		BookService:   NewBookService(repo),
		MemberService: NewMemberService(repo),
	}
}

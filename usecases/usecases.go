// Package usecases contains services and pipelines.
package usecases

import (
	"book-management-system/repositories"
	"book-management-system/usecases/services"
)

// UseCase contains usecases
type UseCase struct {
	Service *services.Services
}

// Init returns UseCase
func Init(repo *repositories.Repository) *UseCase {
	return &UseCase{
		Service: services.Init(repo),
	}
}

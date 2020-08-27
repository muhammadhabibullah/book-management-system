package usecases

import (
	"book-management-system/repositories"
	"book-management-system/usecases/services"
)

// UseCase contains use-cases
type UseCase struct {
	Service *services.Services
}

// Init return UseCase
func Init(repo *repositories.Repository) *UseCase {
	return &UseCase{
		Service: services.Init(repo),
	}
}

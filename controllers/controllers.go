package controllers

import (
	REST "book-management-system/controllers/rest"
	"book-management-system/usecases"
)

// Init set controllers
func Init(useCase *usecases.UseCase) {
	REST.Init(useCase)
}

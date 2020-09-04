// Package controllers contains all controllers
package controllers

import (
	REST "book-management-system/controllers/rest"
	"book-management-system/usecases"
)

// Init sets controllers
func Init(useCase *usecases.UseCase) {
	REST.Init(useCase)
}

func testUnusedFunc() {}

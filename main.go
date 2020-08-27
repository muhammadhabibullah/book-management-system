package main

import (
	"book-management-system/configs"
	"book-management-system/controllers"
	"book-management-system/repositories"
	"book-management-system/usecases"
)

func main() {
	configs.GetConfig()
	controllers.Init(usecases.Init(repositories.Init()))
}

package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"book-management-system/entities/models"
	"book-management-system/usecases"
	"book-management-system/usecases/services"
)

// BookController will handle book domain requests
type BookController struct {
	bookService services.BookService
}

// NewBookController returns new BookController
func NewBookController(route *mux.Router, useCase *usecases.UseCase) {
	ctrl := &BookController{
		bookService: useCase.Service.BookService,
	}

	v1Route := route.PathPrefix("/v1").Subrouter()
	v1BookRoute := v1Route.PathPrefix("/book").Subrouter()
	v1BookRoute.HandleFunc("", ctrl.CreateBook).Methods("POST")
	v1BookRoute.HandleFunc("", ctrl.GetBooks).Methods("GET")
}

// CreateBook handle create book request
func (ctrl *BookController) CreateBook(res http.ResponseWriter, req *http.Request) {
	var book models.Book
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&book); err != nil {
		respondWithError(res, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer func() {
		_ = req.Body.Close()
	}()

	if err := ctrl.bookService.CreateBook(&book); err != nil {
		respondWithError(res, http.StatusBadRequest,
			fmt.Sprintf("Failed create book: %s", err.Error()))
		return
	}

	respondWithJSON(res, http.StatusCreated, book)
}

// GetBooks handle get all books request
func (ctrl *BookController) GetBooks(res http.ResponseWriter, _ *http.Request) {
	books, err := ctrl.bookService.GetBooks()
	if err != nil {
		respondWithError(res, http.StatusBadRequest,
			fmt.Sprintf("Failed create book: %s", err.Error()))
		return
	}

	respondWithJSON(res, http.StatusCreated, books)
}

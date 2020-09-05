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
// @Summary Create a new book
// @Description Create a new book
// @Tags Book
// @Accept json
// @Produce json
// @Param request body models.Book true "Request Body"
// @Success 201 {object} models.Book "Created"
// @Failure 400 {object} responses.ErrorResponse "Bad Request"
// @Failure 500 {object} responses.ErrorResponse "Internal Server Error"
// @Router /v1/book [post]
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
		respondWithError(res, http.StatusInternalServerError,
			fmt.Sprintf("Failed create book: %s", err.Error()))
		return
	}

	respondWithJSON(res, http.StatusCreated, book)
}

// GetBooks handle get all books request
// @Summary Get all books
// @Description Get all books
// @Tags Book
// @Accept json
// @Produce json
// @Success 200 {object} models.Books "OK"
// @Failure 500 {object} responses.ErrorResponse "Internal Server Error"
// @Router /v1/book [get]
func (ctrl *BookController) GetBooks(res http.ResponseWriter, _ *http.Request) {
	books, err := ctrl.bookService.GetBooks()
	if err != nil {
		respondWithError(res, http.StatusInternalServerError,
			fmt.Sprintf("Failed get books: %s", err.Error()))
		return
	}

	respondWithJSON(res, http.StatusOK, books)
}

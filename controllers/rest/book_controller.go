package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"book-management-system/controllers/rest/middlewares"
	"book-management-system/entities/models"
	"book-management-system/usecases"
	"book-management-system/usecases/services"
)

// BookController will handle book domain requests
type BookController struct {
	bookService services.BookService
}

// NewBookController returns new BookController
func NewBookController(
	route *mux.Router,
	useCase *usecases.UseCase,
	m *middlewares.Middleware,
) *BookController {
	ctrl := &BookController{
		bookService: useCase.Service.BookService,
	}

	v1Route := route.PathPrefix("/v1").Subrouter()

	v1BookRoute := v1Route.PathPrefix("/book").Subrouter()

	v1BookAuthorizedRoute := v1BookRoute.PathPrefix("").Subrouter()
	v1BookAuthorizedRoute.Use(m.AuthMiddleware)

	v1BookAuthorizedRoute.HandleFunc("", ctrl.CreateBook).Methods(http.MethodPost)
	v1BookRoute.HandleFunc("", ctrl.GetBooks).Methods(http.MethodGet)
	v1BookAuthorizedRoute.HandleFunc("", ctrl.UpdateBook).Methods(http.MethodPut)

	return ctrl
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
func (ctrl *BookController) CreateBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&book); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := ctrl.bookService.CreateBook(r.Context(), &book); err != nil {
		respondWithError(w, http.StatusInternalServerError,
			fmt.Sprintf("Failed create book: %s", err.Error()))
		return
	}

	respondWithJSON(w, http.StatusCreated, book)
}

// GetBooks handle get all books request
// @Summary Get all books
// @Description Get all books
// @Tags Book
// @Accept json
// @Produce json
// @Param search query string false "Search"
// @Success 200 {object} models.Books "OK"
// @Failure 500 {object} responses.ErrorResponse "Internal Server Error"
// @Router /v1/book [get]
func (ctrl *BookController) GetBooks(w http.ResponseWriter, r *http.Request) {
	var books models.Books
	var err error
	if keyword := r.URL.Query().Get("search"); keyword != "" {
		books, err = ctrl.bookService.SearchBooks(r.Context(), keyword)
	} else {
		books, err = ctrl.bookService.GetBooks(r.Context())
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
			fmt.Sprintf("Failed get books: %s", err.Error()))
		return
	}

	respondWithJSON(w, http.StatusOK, books)
}

// UpdateBook handle update book request
// @Summary Update a book
// @Description Update a book
// @Tags Book
// @Accept json
// @Produce json
// @Param request body models.Book true "Request Body"
// @Success 200 {object} models.Book "Updated"
// @Failure 400 {object} responses.ErrorResponse "Bad Request"
// @Failure 500 {object} responses.ErrorResponse "Internal Server Error"
// @Router /v1/book [put]
func (ctrl *BookController) UpdateBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&book); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := ctrl.bookService.UpdateBook(r.Context(), &book); err != nil {
		respondWithError(w, http.StatusInternalServerError,
			fmt.Sprintf("Failed update book: %s", err.Error()))
		return
	}

	respondWithJSON(w, http.StatusOK, book)
}

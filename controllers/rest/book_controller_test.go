package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"

	"book-management-system/controllers/rest/middlewares"
	"book-management-system/controllers/rest/responses"
	"book-management-system/entities/models"
	mocks "book-management-system/mocks/services"
	"book-management-system/repositories"
	"book-management-system/usecases"
	"book-management-system/usecases/services"
)

func TestNewBookController(t *testing.T) {
	repo := &repositories.Repository{}
	bookService := services.NewBookService(repo)
	usecase := &usecases.UseCase{
		Service: &services.Services{
			BookService: bookService,
		},
	}

	route := mux.NewRouter()
	m := &middlewares.Middleware{}
	got := NewBookController(route, usecase, m)
	expected := &BookController{
		bookService: bookService,
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("NewBookController returns %+v\n expected %+v",
			got, expected)
	}
}

const (
	v1BookURL = "/v1/book"
)

func TestBookControllerCreateBook(t *testing.T) {
	type input struct {
		valid              bool
		ctx                context.Context
		requestBody        *models.Book
		invalidRequestBody models.Books
	}
	type output struct {
		responseBody interface{}
	}
	type mockConfig struct {
		given input
		mock  *mocks.MockBookService
	}

	tests := []struct {
		name           string
		givenInput     input
		expectedOutput output
		configureMock  func(mockConfig)
	}{
		{
			name: "failed: invalid request body",
			givenInput: input{
				valid: false,
				ctx:   context.TODO(),
				invalidRequestBody: models.Books{
					{
						Name: "C++",
						ISBN: "1234",
					},
				},
			},
			expectedOutput: output{
				responseBody: responses.ErrorResponse{
					"error": "Invalid request payload",
				},
			},
			configureMock: func(mockConfig) {
				// do nothing
			},
		},
		{
			name: "failed: create book service returns error",
			givenInput: input{
				valid: true,
				ctx:   context.TODO(),
				requestBody: &models.Book{
					Name: "C++",
					ISBN: "1234",
				},
			},
			expectedOutput: output{
				responseBody: responses.ErrorResponse{
					"error": fmt.Sprintf("Failed create book: %s", errService.Error()),
				},
			},
			configureMock: func(conf mockConfig) {
				conf.mock.EXPECT().
					CreateBook(conf.given.ctx, conf.given.requestBody).
					Return(errService)
			},
		},
		{
			name: "success: create book",
			givenInput: input{
				valid: true,
				ctx:   context.TODO(),
				requestBody: &models.Book{
					Name: "C++",
					ISBN: "1234",
				},
			},
			expectedOutput: output{
				responseBody: models.Book{
					Name: "C++",
					ISBN: "1234",
				},
			},
			configureMock: func(conf mockConfig) {
				conf.mock.EXPECT().
					CreateBook(conf.given.ctx, conf.given.requestBody).
					Return(nil)
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var marshalledRequestBody []byte
			if tt.givenInput.valid {
				marshalledRequestBody, _ = json.Marshal(tt.givenInput.requestBody)
			} else {
				marshalledRequestBody, _ = json.Marshal(tt.givenInput.invalidRequestBody)
			}

			req, _ := http.NewRequest(
				http.MethodPost,
				v1BookURL,
				bytes.NewBuffer(marshalledRequestBody),
			)
			resp := httptest.NewRecorder()

			bookServiceMock := mocks.NewMockBookService(ctrl)
			tt.configureMock(mockConfig{
				given: tt.givenInput,
				mock:  bookServiceMock,
			})

			bookController := &BookController{
				bookService: bookServiceMock,
			}

			bookController.CreateBook(resp, req)

			got := resp.Body.String()
			expected, _ := json.Marshal(tt.expectedOutput.responseBody)
			if got != string(expected) {
				t.Errorf("CreateBook() got response body %s\n expected %s",
					got, string(expected))
			}
		})
	}
}

func TestBookControllerGetBook(t *testing.T) {
	type input struct {
		ctx            context.Context
		httpRequestURL string
		query          string
	}
	type output struct {
		responseBody interface{}
	}
	type mockConfig struct {
		given    input
		expected output
		mock     *mocks.MockBookService
	}

	tests := []struct {
		name           string
		givenInput     input
		expectedOutput output
		configureMock  func(mockConfig)
	}{
		{
			name: "failed: get books service returns error",
			givenInput: input{
				ctx:            context.TODO(),
				httpRequestURL: v1BookURL,
			},
			expectedOutput: output{
				responseBody: responses.ErrorResponse{
					"error": fmt.Sprintf("Failed get books: %s", errService.Error()),
				},
			},
			configureMock: func(conf mockConfig) {
				conf.mock.EXPECT().
					GetBooks(conf.given.ctx).
					Return(models.Books{}, errService)
			},
		},
		{
			name: "success: get books",
			givenInput: input{
				ctx:            context.TODO(),
				httpRequestURL: v1BookURL,
			},
			expectedOutput: output{
				responseBody: models.Books{
					{
						Name: "C++",
						ISBN: "1234",
					},
				},
			},
			configureMock: func(conf mockConfig) {
				conf.mock.EXPECT().
					GetBooks(conf.given.ctx).
					Return(conf.expected.responseBody, nil)
			},
		},
		{
			name: "success: search books",
			givenInput: input{
				ctx:            context.TODO(),
				httpRequestURL: v1BookURL + "?search=1234",
				query:          "1234",
			},
			expectedOutput: output{
				responseBody: models.Books{
					{
						Name: "C++",
						ISBN: "1234",
					},
				},
			},
			configureMock: func(conf mockConfig) {
				conf.mock.EXPECT().
					SearchBooks(conf.given.ctx, conf.given.query).
					Return(conf.expected.responseBody, nil)
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(
				http.MethodGet,
				tt.givenInput.httpRequestURL,
				nil,
			)
			resp := httptest.NewRecorder()

			bookServiceMock := mocks.NewMockBookService(ctrl)
			tt.configureMock(mockConfig{
				given:    tt.givenInput,
				expected: tt.expectedOutput,
				mock:     bookServiceMock,
			})

			bookController := &BookController{
				bookService: bookServiceMock,
			}

			bookController.GetBooks(resp, req)

			got := resp.Body.String()
			expected, _ := json.Marshal(tt.expectedOutput.responseBody)
			if got != string(expected) {
				t.Errorf("GetBooks() got response body %s\n expected %s",
					got, string(expected))
			}
		})
	}
}

func TestBookControllerUpdateBook(t *testing.T) {
	type input struct {
		valid              bool
		ctx                context.Context
		requestBody        *models.Book
		invalidRequestBody models.Books
	}
	type output struct {
		responseBody interface{}
	}
	type confMock struct {
		given input
		mock  *mocks.MockBookService
	}

	tests := []struct {
		name           string
		givenInput     input
		expectedOutput output
		configureMock  func(confMock)
	}{
		{
			name: "failed: invalid request body",
			givenInput: input{
				valid: false,
				ctx:   context.TODO(),
				invalidRequestBody: models.Books{
					{
						Name: "C++",
						ISBN: "1234",
					},
				},
			},
			expectedOutput: output{
				responseBody: responses.ErrorResponse{
					"error": "Invalid request payload",
				},
			},
			configureMock: func(confMock) {
				// do nothing
			},
		},
		{
			name: "failed: update book service returns error",
			givenInput: input{
				valid: true,
				ctx:   context.TODO(),
				requestBody: &models.Book{
					Name: "C++",
					ISBN: "1234",
				},
			},
			expectedOutput: output{
				responseBody: responses.ErrorResponse{
					"error": fmt.Sprintf("Failed update book: %s", errService.Error()),
				},
			},
			configureMock: func(conf confMock) {
				conf.mock.EXPECT().
					UpdateBook(conf.given.ctx, conf.given.requestBody).
					Return(errService)
			},
		},
		{
			name: "success: update book",
			givenInput: input{
				valid: true,
				ctx:   context.TODO(),
				requestBody: &models.Book{
					Name: "C++",
					ISBN: "1234",
				},
			},
			expectedOutput: output{
				responseBody: models.Book{
					Name: "C++",
					ISBN: "1234",
				},
			},
			configureMock: func(conf confMock) {
				conf.mock.EXPECT().
					UpdateBook(conf.given.ctx, conf.given.requestBody).
					Return(nil)
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var marshalledRequestBody []byte
			if tt.givenInput.valid {
				marshalledRequestBody, _ = json.Marshal(tt.givenInput.requestBody)
			} else {
				marshalledRequestBody, _ = json.Marshal(tt.givenInput.invalidRequestBody)
			}

			req, _ := http.NewRequest(
				http.MethodPut,
				v1BookURL,
				bytes.NewBuffer(marshalledRequestBody),
			)
			resp := httptest.NewRecorder()

			bookServiceMock := mocks.NewMockBookService(ctrl)
			tt.configureMock(confMock{
				given: tt.givenInput,
				mock:  bookServiceMock,
			})

			bookController := &BookController{
				bookService: bookServiceMock,
			}

			bookController.UpdateBook(resp, req)

			got := resp.Body.String()
			expected, _ := json.Marshal(tt.expectedOutput.responseBody)
			if got != string(expected) {
				t.Errorf("UpdateBook() got response body %s\n expected %s",
					got, string(expected))
			}
		})
	}
}

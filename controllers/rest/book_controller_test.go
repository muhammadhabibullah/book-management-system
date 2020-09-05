package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"

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
	got := NewBookController(route, usecase)
	expected := &BookController{
		bookService: bookService,
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("NewBookController returns %+v\n expected %+v",
			got, expected)
	}
}

func TestBookController_CreateBook(t *testing.T) {
	type input struct {
		valid              bool
		requestBody        *models.Book
		invalidRequestBody models.Books
	}
	type output struct {
		responseBody interface{}
	}
	type confMock struct {
		input input
		mock  *mocks.MockBookService
	}

	tests := []struct {
		name          string
		input         input
		output        output
		configureMock func(confMock)
	}{
		{
			name: "failed: invalid request body",
			input: input{
				valid: false,
				invalidRequestBody: models.Books{
					{
						Name: "C++",
						ISBN: "1234",
					},
				},
			},
			output: output{
				responseBody: responses.ErrorResponse{
					"error": "Invalid request payload",
				},
			},
			configureMock: func(confMock) {},
		},
		{
			name: "failed: service error",
			input: input{
				valid: true,
				requestBody: &models.Book{
					Name: "C++",
					ISBN: "1234",
				},
			},
			output: output{
				responseBody: responses.ErrorResponse{
					"error": fmt.Sprintf("Failed create book: %s", "service error"),
				},
			},
			configureMock: func(conf confMock) {
				conf.mock.EXPECT().
					CreateBook(conf.input.requestBody).
					Return(errors.New("service error"))
			},
		},
		{
			name: "success: create book",
			input: input{
				valid: true,
				requestBody: &models.Book{
					Name: "C++",
					ISBN: "1234",
				},
			},
			output: output{
				responseBody: models.Book{
					Name: "C++",
					ISBN: "1234",
				},
			},
			configureMock: func(conf confMock) {
				conf.mock.EXPECT().
					CreateBook(conf.input.requestBody).
					Return(nil)
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var marshalledRequestBody []byte
			if tt.input.valid {
				marshalledRequestBody, _ = json.Marshal(tt.input.requestBody)
			} else {
				marshalledRequestBody, _ = json.Marshal(tt.input.invalidRequestBody)
			}

			req, _ := http.NewRequest(
				"POST",
				"/v1/book",
				bytes.NewBuffer(marshalledRequestBody),
			)
			resp := httptest.NewRecorder()

			bookServiceMock := mocks.NewMockBookService(ctrl)
			tt.configureMock(confMock{
				input: tt.input,
				mock:  bookServiceMock,
			})

			bookController := &BookController{
				bookService: bookServiceMock,
			}

			bookController.CreateBook(resp, req)

			got := resp.Body.String()
			expected, _ := json.Marshal(tt.output.responseBody)
			if got != string(expected) {
				t.Errorf("got response body %s\n expected %s",
					got, string(expected))
			}
		})
	}
}

func TestBookController_GetBook(t *testing.T) {
	type output struct {
		responseBody interface{}
	}
	type confMock struct {
		mock *mocks.MockBookService
	}

	tests := []struct {
		name          string
		output        output
		configureMock func(confMock)
	}{
		{
			name: "failed: service error",
			output: output{
				responseBody: responses.ErrorResponse{
					"error": fmt.Sprintf("Failed get books: %s", "service error"),
				},
			},
			configureMock: func(conf confMock) {
				conf.mock.EXPECT().
					GetBooks().
					Return(models.Books{}, errors.New("service error"))
			},
		},
		{
			name: "success: create book",
			output: output{
				responseBody: models.Books{
					{
						Name: "C++",
						ISBN: "1234",
					},
				},
			},
			configureMock: func(conf confMock) {
				conf.mock.EXPECT().
					GetBooks().
					Return(models.Books{
						{
							Name: "C++",
							ISBN: "1234",
						},
					}, nil)
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(
				"GET",
				"/v1/book",
				nil,
			)
			resp := httptest.NewRecorder()

			bookServiceMock := mocks.NewMockBookService(ctrl)
			tt.configureMock(confMock{
				mock: bookServiceMock,
			})

			bookController := &BookController{
				bookService: bookServiceMock,
			}

			bookController.GetBooks(resp, req)

			got := resp.Body.String()
			expected, _ := json.Marshal(tt.output.responseBody)
			if got != string(expected) {
				t.Errorf("got response body %s\n expected %s",
					got, string(expected))
			}
		})
	}
}

func TestBookController_UpdateBook(t *testing.T) {
	type input struct {
		valid              bool
		requestBody        *models.Book
		invalidRequestBody models.Books
	}
	type output struct {
		responseBody interface{}
	}
	type confMock struct {
		input input
		mock  *mocks.MockBookService
	}

	tests := []struct {
		name          string
		input         input
		output        output
		configureMock func(confMock)
	}{
		{
			name: "failed: invalid request body",
			input: input{
				valid: false,
				invalidRequestBody: models.Books{
					{
						Name: "C++",
						ISBN: "1234",
					},
				},
			},
			output: output{
				responseBody: responses.ErrorResponse{
					"error": "Invalid request payload",
				},
			},
			configureMock: func(confMock) {},
		},
		{
			name: "failed: service error",
			input: input{
				valid: true,
				requestBody: &models.Book{
					Name: "C++",
					ISBN: "1234",
				},
			},
			output: output{
				responseBody: responses.ErrorResponse{
					"error": fmt.Sprintf("Failed update book: %s", "service error"),
				},
			},
			configureMock: func(conf confMock) {
				conf.mock.EXPECT().
					UpdateBook(conf.input.requestBody).
					Return(errors.New("service error"))
			},
		},
		{
			name: "success: update book",
			input: input{
				valid: true,
				requestBody: &models.Book{
					Name: "C++",
					ISBN: "1234",
				},
			},
			output: output{
				responseBody: models.Book{
					Name: "C++",
					ISBN: "1234",
				},
			},
			configureMock: func(conf confMock) {
				conf.mock.EXPECT().
					UpdateBook(conf.input.requestBody).
					Return(nil)
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var marshalledRequestBody []byte
			if tt.input.valid {
				marshalledRequestBody, _ = json.Marshal(tt.input.requestBody)
			} else {
				marshalledRequestBody, _ = json.Marshal(tt.input.invalidRequestBody)
			}

			req, _ := http.NewRequest(
				"PUT",
				"/v1/book",
				bytes.NewBuffer(marshalledRequestBody),
			)
			resp := httptest.NewRecorder()

			bookServiceMock := mocks.NewMockBookService(ctrl)
			tt.configureMock(confMock{
				input: tt.input,
				mock:  bookServiceMock,
			})

			bookController := &BookController{
				bookService: bookServiceMock,
			}

			bookController.UpdateBook(resp, req)

			got := resp.Body.String()
			expected, _ := json.Marshal(tt.output.responseBody)
			if got != string(expected) {
				t.Errorf("got response body %s\n expected %s",
					got, string(expected))
			}
		})
	}
}

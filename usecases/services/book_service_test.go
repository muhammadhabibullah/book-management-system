package services

import (
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"book-management-system/entities/models"
	mocks "book-management-system/mocks/repositories"
	"book-management-system/repositories"
	"book-management-system/repositories/mysql"
)

func TestNewBookService(t *testing.T) {
	bookRepo := mysql.NewBookRepository(nil)
	repo := &repositories.Repository{
		BookRepository: bookRepo,
	}

	got := NewBookService(repo)
	expected := &bookService{
		BookRepository: bookRepo,
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("NewBookService returns %+v\n expected %+v",
			got, expected)
	}
	if _, ok := got.(BookService); !ok {
		t.Errorf("NewBookService returns object not implements BookService")
	}
}

var errDatabase = errors.New("error database")

func TestBookService_CreateBook(t *testing.T) {
	type input struct {
		book *models.Book
	}
	type output struct {
		err error
	}
	type confMock struct {
		input        input
		output       output
		bookRepoMock *mocks.MockBookRepository
	}

	tests := []struct {
		name          string
		input         input
		output        output
		configureMock func(confMock)
	}{
		{
			name: "success create book",
			input: input{
				book: &models.Book{
					Name: "C++",
					ISBN: "1234",
				},
			},
			output: output{
				err: nil,
			},
			configureMock: func(conf confMock) {
				conf.bookRepoMock.EXPECT().
					CreateBook(conf.input.book).
					Return(conf.output.err)
			},
		},
		{
			name: "failed create book",
			input: input{
				book: &models.Book{
					Name: "C++",
					ISBN: "1234",
				},
			},
			output: output{
				err: errDatabase,
			},
			configureMock: func(conf confMock) {
				conf.bookRepoMock.EXPECT().
					CreateBook(conf.input.book).
					Return(conf.output.err)
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bookRepoMock := mocks.NewMockBookRepository(ctrl)

			bookService := NewBookService(
				&repositories.Repository{
					BookRepository: bookRepoMock,
				},
			)

			tt.configureMock(confMock{
				input:        tt.input,
				output:       tt.output,
				bookRepoMock: bookRepoMock,
			})

			err := bookService.CreateBook(tt.input.book)
			if !errors.Is(err, tt.output.err) {
				t.Errorf("unexpected error %+v, expected %+v",
					err, tt.output.err)
			}
		})
	}
}

func TestBookService_GetBook(t *testing.T) {
	type output struct {
		books models.Books
		err   error
	}
	type confMock struct {
		output       output
		bookRepoMock *mocks.MockBookRepository
	}

	tests := []struct {
		name          string
		output        output
		configureMock func(confMock)
	}{
		{
			name: "success get book",
			output: output{
				books: models.Books{
					{
						Name: "C++",
						ISBN: "1234",
					},
				},
				err: nil,
			},
			configureMock: func(conf confMock) {
				conf.bookRepoMock.EXPECT().
					GetAll().
					Return(
						conf.output.books,
						conf.output.err,
					)
			},
		},
		{
			name: "failed get book",
			output: output{
				books: models.Books{
					{
						Name: "C++",
						ISBN: "1234",
					},
				},
				err: errDatabase,
			},
			configureMock: func(conf confMock) {
				conf.bookRepoMock.EXPECT().
					GetAll().
					Return(
						conf.output.books,
						conf.output.err,
					)
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bookRepoMock := mocks.NewMockBookRepository(ctrl)

			bookService := NewBookService(
				&repositories.Repository{
					BookRepository: bookRepoMock,
				},
			)

			tt.configureMock(confMock{
				output:       tt.output,
				bookRepoMock: bookRepoMock,
			})

			got, err := bookService.GetBooks()
			if !errors.Is(err, tt.output.err) {
				t.Errorf("unexpected error %+v, expected %+v",
					err, tt.output.err)
			}
			if !reflect.DeepEqual(got, tt.output.books) {
				t.Errorf("result got %+v, expected %+v",
					got, tt.output.books)
			}
		})
	}
}

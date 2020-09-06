package services

import (
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"book-management-system/entities/models"
	esMocks "book-management-system/mocks/repositories/elasticsearch"
	mySqlMocks "book-management-system/mocks/repositories/mysql"
	"book-management-system/repositories"
	"book-management-system/repositories/elasticsearch"
	"book-management-system/repositories/mysql"
)

func TestNewBookService(t *testing.T) {
	mySQLBookRepo := mysql.NewBookRepository(nil)
	esBookRepo := elasticsearch.NewBookRepository(nil)
	repo := &repositories.Repository{
		MySQLBookRepository: mySQLBookRepo,
		ESBookRepository:    esBookRepo,
	}

	got := NewBookService(repo)
	expected := &bookService{
		MySQLBookRepository: mySQLBookRepo,
		ESBookRepository:    esBookRepo,
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
		input             input
		output            output
		mySQLBookRepoMock *mySqlMocks.MockBookRepository
		esBookRepoMock    *esMocks.MockBookRepository
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
				conf.mySQLBookRepoMock.EXPECT().
					CreateBook(conf.input.book).
					Return(conf.output.err)

				conf.esBookRepoMock.EXPECT().
					IndexBook(conf.input.book).
					Return(nil)
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
				conf.mySQLBookRepoMock.EXPECT().
					CreateBook(conf.input.book).
					Return(conf.output.err)
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mySQLBookRepoMock := mySqlMocks.NewMockBookRepository(ctrl)
			esBookRepoMock := esMocks.NewMockBookRepository(ctrl)

			bookService := NewBookService(
				&repositories.Repository{
					MySQLBookRepository: mySQLBookRepoMock,
					ESBookRepository:    esBookRepoMock,
				},
			)

			tt.configureMock(confMock{
				input:             tt.input,
				output:            tt.output,
				mySQLBookRepoMock: mySQLBookRepoMock,
				esBookRepoMock:    esBookRepoMock,
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
		output            output
		mySQLBookRepoMock *mySqlMocks.MockBookRepository
	}

	tests := []struct {
		name          string
		output        output
		configureMock func(confMock)
	}{
		{
			name: "success get books",
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
				conf.mySQLBookRepoMock.EXPECT().
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
				books: models.Books{},
				err:   errDatabase,
			},
			configureMock: func(conf confMock) {
				conf.mySQLBookRepoMock.EXPECT().
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
			mySQLBookRepoMock := mySqlMocks.NewMockBookRepository(ctrl)

			bookService := NewBookService(
				&repositories.Repository{
					MySQLBookRepository: mySQLBookRepoMock,
				},
			)

			tt.configureMock(confMock{
				output:            tt.output,
				mySQLBookRepoMock: mySQLBookRepoMock,
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

func TestBookService_UpdateBook(t *testing.T) {
	type input struct {
		book *models.Book
	}
	type output struct {
		err error
	}
	type confMock struct {
		input             input
		output            output
		mySQLBookRepoMock *mySqlMocks.MockBookRepository
		esBookRepoMock    *esMocks.MockBookRepository
	}

	tests := []struct {
		name          string
		input         input
		output        output
		configureMock func(confMock)
	}{
		{
			name: "success update book",
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
				conf.mySQLBookRepoMock.EXPECT().
					UpdateBook(conf.input.book).
					Return(conf.output.err)

				conf.esBookRepoMock.EXPECT().
					IndexBook(conf.input.book).
					Return(nil)
			},
		},
		{
			name: "failed update book",
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
				conf.mySQLBookRepoMock.EXPECT().
					UpdateBook(conf.input.book).
					Return(conf.output.err)
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mySQLBookRepoMock := mySqlMocks.NewMockBookRepository(ctrl)
			esBookRepoMock := esMocks.NewMockBookRepository(ctrl)

			bookService := NewBookService(
				&repositories.Repository{
					MySQLBookRepository: mySQLBookRepoMock,
					ESBookRepository:    esBookRepoMock,
				},
			)

			tt.configureMock(confMock{
				input:             tt.input,
				output:            tt.output,
				mySQLBookRepoMock: mySQLBookRepoMock,
				esBookRepoMock:    esBookRepoMock,
			})

			err := bookService.UpdateBook(tt.input.book)
			if !errors.Is(err, tt.output.err) {
				t.Errorf("unexpected error %+v, expected %+v",
					err, tt.output.err)
			}
		})
	}
}

func TestBookService_SearchBook(t *testing.T) {
	type input struct {
		keyword string
	}
	type output struct {
		books models.Books
		err   error
	}
	type confMock struct {
		input          input
		output         output
		esBookRepoMock *esMocks.MockBookRepository
	}

	tests := []struct {
		name          string
		input         input
		output        output
		configureMock func(confMock)
	}{
		{
			name: "success search books",
			input: input{
				keyword: "C++",
			},
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
				conf.esBookRepoMock.EXPECT().
					SearchBook(conf.input.keyword).
					Return(
						conf.output.books,
						conf.output.err,
					)
			},
		},
		{
			name: "failed search books",
			input: input{
				keyword: "C++",
			},
			output: output{
				books: models.Books{},
				err:   errDatabase,
			},
			configureMock: func(conf confMock) {
				conf.esBookRepoMock.EXPECT().
					SearchBook(conf.input.keyword).
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
			esBookRepoMock := esMocks.NewMockBookRepository(ctrl)

			bookService := NewBookService(
				&repositories.Repository{
					ESBookRepository: esBookRepoMock,
				},
			)

			tt.configureMock(confMock{
				input:          tt.input,
				output:         tt.output,
				esBookRepoMock: esBookRepoMock,
			})

			got, err := bookService.SearchBooks(tt.input.keyword)
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

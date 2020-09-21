package services

import (
	"errors"
	"reflect"
	"sync"
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
	type mockConfig struct {
		wg                *sync.WaitGroup
		given             input
		expected          output
		mySQLBookRepoMock *mySqlMocks.MockBookRepository
		esBookRepoMock    *esMocks.MockBookRepository
	}

	tests := []struct {
		name           string
		givenInput     input
		expectedOutput output
		configureMock  func(mockConfig)
	}{
		{
			name: "success create book",
			givenInput: input{
				book: &models.Book{
					Name: "C++",
					ISBN: "1234",
				},
			},
			expectedOutput: output{
				err: nil,
			},
			configureMock: func(conf mockConfig) {
				conf.mySQLBookRepoMock.EXPECT().
					CreateBook(conf.given.book).
					Return(conf.expected.err)

				conf.wg.Add(1)

				conf.esBookRepoMock.EXPECT().
					IndexBook(conf.given.book).
					DoAndReturn(func(*models.Book) error {
						conf.wg.Done()
						return nil
					})
			},
		},
		{
			name: "failed create book",
			givenInput: input{
				book: &models.Book{
					Name: "C++",
					ISBN: "1234",
				},
			},
			expectedOutput: output{
				err: errDatabase,
			},
			configureMock: func(conf mockConfig) {
				conf.mySQLBookRepoMock.EXPECT().
					CreateBook(conf.given.book).
					Return(conf.expected.err)
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
			wg := sync.WaitGroup{}

			tt.configureMock(mockConfig{
				wg:                &wg,
				given:             tt.givenInput,
				expected:          tt.expectedOutput,
				mySQLBookRepoMock: mySQLBookRepoMock,
				esBookRepoMock:    esBookRepoMock,
			})

			err := bookService.CreateBook(tt.givenInput.book)
			wg.Wait()
			if expectedError := tt.expectedOutput.err; !errors.Is(err, expectedError) {
				t.Errorf("unexpected error %+v, expected %+v",
					err, expectedError)
			}
		})
	}
}

func TestBookService_GetBook(t *testing.T) {
	type output struct {
		books models.Books
		err   error
	}
	type mockConfig struct {
		expected          output
		mySQLBookRepoMock *mySqlMocks.MockBookRepository
	}

	tests := []struct {
		name           string
		expectedOutput output
		configureMock  func(mockConfig)
	}{
		{
			name: "success get books",
			expectedOutput: output{
				books: models.Books{
					{
						Name: "C++",
						ISBN: "1234",
					},
				},
				err: nil,
			},
			configureMock: func(conf mockConfig) {
				conf.mySQLBookRepoMock.EXPECT().
					GetAll().
					Return(
						conf.expected.books,
						conf.expected.err,
					)
			},
		},
		{
			name: "failed get book",
			expectedOutput: output{
				books: models.Books{},
				err:   errDatabase,
			},
			configureMock: func(conf mockConfig) {
				conf.mySQLBookRepoMock.EXPECT().
					GetAll().
					Return(
						conf.expected.books,
						conf.expected.err,
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

			tt.configureMock(mockConfig{
				expected:          tt.expectedOutput,
				mySQLBookRepoMock: mySQLBookRepoMock,
			})

			books, err := bookService.GetBooks()
			if expectedError := tt.expectedOutput.err; !errors.Is(err, expectedError) {
				t.Errorf("unexpected error %+v, expected %+v",
					err, expectedError)
			}
			if expectedBooks := tt.expectedOutput.books; !reflect.DeepEqual(books, expectedBooks) {
				t.Errorf("got books %+v, expected %+v",
					books, expectedBooks)
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
	type mockConfig struct {
		wg                *sync.WaitGroup
		given             input
		expected          output
		mySQLBookRepoMock *mySqlMocks.MockBookRepository
		esBookRepoMock    *esMocks.MockBookRepository
	}

	tests := []struct {
		name           string
		givenInput     input
		expectedOutput output
		configureMock  func(mockConfig)
	}{
		{
			name: "success update book",
			givenInput: input{
				book: &models.Book{
					Name: "C++",
					ISBN: "1234",
				},
			},
			expectedOutput: output{
				err: nil,
			},
			configureMock: func(conf mockConfig) {
				conf.mySQLBookRepoMock.EXPECT().
					UpdateBook(conf.given.book).
					Return(conf.expected.err)

				conf.wg.Add(1)

				conf.esBookRepoMock.EXPECT().
					IndexBook(conf.given.book).
					DoAndReturn(func(*models.Book) error {
						conf.wg.Done()
						return nil
					})
			},
		},
		{
			name: "failed update book",
			givenInput: input{
				book: &models.Book{
					Name: "C++",
					ISBN: "1234",
				},
			},
			expectedOutput: output{
				err: errDatabase,
			},
			configureMock: func(conf mockConfig) {
				conf.mySQLBookRepoMock.EXPECT().
					UpdateBook(conf.given.book).
					Return(conf.expected.err)
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
			wg := sync.WaitGroup{}

			tt.configureMock(mockConfig{
				wg:                &wg,
				given:             tt.givenInput,
				expected:          tt.expectedOutput,
				mySQLBookRepoMock: mySQLBookRepoMock,
				esBookRepoMock:    esBookRepoMock,
			})

			err := bookService.UpdateBook(tt.givenInput.book)
			wg.Wait()
			if expectedError := tt.expectedOutput.err; !errors.Is(err, expectedError) {
				t.Errorf("unexpected error %+v, expected %+v",
					err, expectedError)
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
	type mockConfig struct {
		given          input
		expected       output
		esBookRepoMock *esMocks.MockBookRepository
	}

	tests := []struct {
		name           string
		givenInput     input
		expectedOutput output
		configureMock  func(mockConfig)
	}{
		{
			name: "success search books",
			givenInput: input{
				keyword: "C++",
			},
			expectedOutput: output{
				books: models.Books{
					{
						Name: "C++",
						ISBN: "1234",
					},
				},
				err: nil,
			},
			configureMock: func(conf mockConfig) {
				conf.esBookRepoMock.EXPECT().
					SearchBook(conf.given.keyword).
					Return(
						conf.expected.books,
						conf.expected.err,
					)
			},
		},
		{
			name: "failed search books",
			givenInput: input{
				keyword: "C++",
			},
			expectedOutput: output{
				books: models.Books{},
				err:   errDatabase,
			},
			configureMock: func(conf mockConfig) {
				conf.esBookRepoMock.EXPECT().
					SearchBook(conf.given.keyword).
					Return(
						conf.expected.books,
						conf.expected.err,
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

			tt.configureMock(mockConfig{
				given:          tt.givenInput,
				expected:       tt.expectedOutput,
				esBookRepoMock: esBookRepoMock,
			})

			books, err := bookService.SearchBooks(tt.givenInput.keyword)
			if expectedError := tt.expectedOutput.err; !errors.Is(err, expectedError) {
				t.Errorf("unexpected error %+v, expected %+v",
					err, expectedError)
			}
			if expectedBooks := tt.expectedOutput.books; !reflect.DeepEqual(books, expectedBooks) {
				t.Errorf("got books %+v, expected %+v",
					books, expectedBooks)
			}
		})
	}
}

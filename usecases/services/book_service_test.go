package services

import (
	"context"
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

func TestBookServiceCreateBook(t *testing.T) {
	type input struct {
		ctx  context.Context
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
				ctx: context.TODO(),
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
					CreateBook(gomock.Any(), conf.given.book).
					Return(conf.expected.err)

				conf.wg.Add(1)

				conf.esBookRepoMock.EXPECT().
					IndexBook(context.Background(), conf.given.book).
					DoAndReturn(func(interface{}, *models.Book) error {
						conf.wg.Done()
						return nil
					})
			},
		},
		{
			name: "failed create book",
			givenInput: input{
				ctx: context.TODO(),
				book: &models.Book{
					Name: "C++",
					ISBN: "1234",
				},
			},
			expectedOutput: output{
				err: errRepository,
			},
			configureMock: func(conf mockConfig) {
				conf.mySQLBookRepoMock.EXPECT().
					CreateBook(gomock.Any(), conf.given.book).
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

			bookService := &bookService{
				MySQLBookRepository: mySQLBookRepoMock,
				ESBookRepository:    esBookRepoMock,
			}
			wg := sync.WaitGroup{}

			tt.configureMock(mockConfig{
				wg:                &wg,
				given:             tt.givenInput,
				expected:          tt.expectedOutput,
				mySQLBookRepoMock: mySQLBookRepoMock,
				esBookRepoMock:    esBookRepoMock,
			})

			err := bookService.CreateBook(tt.givenInput.ctx, tt.givenInput.book)
			wg.Wait()
			if expectedError := tt.expectedOutput.err; !errors.Is(err, expectedError) {
				t.Errorf("CreateBook() got error %+v, expected %+v",
					err, expectedError)
			}
		})
	}
}

func TestBookServiceGetBook(t *testing.T) {
	type input struct {
		ctx context.Context
	}
	type output struct {
		books models.Books
		err   error
	}
	type mockConfig struct {
		given             input
		expected          output
		mySQLBookRepoMock *mySqlMocks.MockBookRepository
	}

	tests := []struct {
		name           string
		givenInput     input
		expectedOutput output
		configureMock  func(mockConfig)
	}{
		{
			name: "success get books",
			givenInput: input{
				ctx: context.TODO(),
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
				conf.mySQLBookRepoMock.EXPECT().
					GetAll(gomock.Any()).
					Return(
						conf.expected.books,
						conf.expected.err,
					)
			},
		},
		{
			name: "failed get book",
			givenInput: input{
				ctx: context.TODO(),
			},
			expectedOutput: output{
				books: models.Books{},
				err:   errRepository,
			},
			configureMock: func(conf mockConfig) {
				conf.mySQLBookRepoMock.EXPECT().
					GetAll(gomock.Any()).
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

			bookService := &bookService{
				MySQLBookRepository: mySQLBookRepoMock,
			}

			tt.configureMock(mockConfig{
				given:             tt.givenInput,
				expected:          tt.expectedOutput,
				mySQLBookRepoMock: mySQLBookRepoMock,
			})

			books, err := bookService.GetBooks(tt.givenInput.ctx)
			if expectedError := tt.expectedOutput.err; !errors.Is(err, expectedError) {
				t.Errorf("GetBooks() got error %+v, expected %+v",
					err, expectedError)
			}
			if expectedBooks := tt.expectedOutput.books; !reflect.DeepEqual(books, expectedBooks) {
				t.Errorf("GetBooks() got books %+v, expected %+v",
					books, expectedBooks)
			}
		})
	}
}

func TestBookServiceUpdateBook(t *testing.T) {
	type input struct {
		ctx  context.Context
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
				ctx: context.TODO(),
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
					UpdateBook(gomock.Any(), conf.given.book).
					Return(conf.expected.err)

				conf.wg.Add(1)

				conf.esBookRepoMock.EXPECT().
					IndexBook(context.Background(), conf.given.book).
					DoAndReturn(func(interface{}, *models.Book) error {
						conf.wg.Done()
						return nil
					})
			},
		},
		{
			name: "failed update book",
			givenInput: input{
				ctx: context.TODO(),
				book: &models.Book{
					Name: "C++",
					ISBN: "1234",
				},
			},
			expectedOutput: output{
				err: errRepository,
			},
			configureMock: func(conf mockConfig) {
				conf.mySQLBookRepoMock.EXPECT().
					UpdateBook(gomock.Any(), conf.given.book).
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

			bookService := &bookService{
				MySQLBookRepository: mySQLBookRepoMock,
				ESBookRepository:    esBookRepoMock,
			}
			wg := sync.WaitGroup{}

			tt.configureMock(mockConfig{
				wg:                &wg,
				given:             tt.givenInput,
				expected:          tt.expectedOutput,
				mySQLBookRepoMock: mySQLBookRepoMock,
				esBookRepoMock:    esBookRepoMock,
			})

			err := bookService.UpdateBook(tt.givenInput.ctx, tt.givenInput.book)
			wg.Wait()
			if expectedError := tt.expectedOutput.err; !errors.Is(err, expectedError) {
				t.Errorf("UpdateBook() got error %+v, expected %+v",
					err, expectedError)
			}
		})
	}
}

func TestBookServiceSearchBook(t *testing.T) {
	type input struct {
		ctx     context.Context
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
				ctx:     context.TODO(),
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
					SearchBook(gomock.Any(), conf.given.keyword).
					Return(
						conf.expected.books,
						conf.expected.err,
					)
			},
		},
		{
			name: "failed search books",
			givenInput: input{
				ctx:     context.TODO(),
				keyword: "C++",
			},
			expectedOutput: output{
				books: models.Books{},
				err:   errRepository,
			},
			configureMock: func(conf mockConfig) {
				conf.esBookRepoMock.EXPECT().
					SearchBook(gomock.Any(), conf.given.keyword).
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

			bookService := &bookService{
				ESBookRepository: esBookRepoMock,
			}

			tt.configureMock(mockConfig{
				given:          tt.givenInput,
				expected:       tt.expectedOutput,
				esBookRepoMock: esBookRepoMock,
			})

			books, err := bookService.SearchBooks(tt.givenInput.ctx, tt.givenInput.keyword)
			if expectedError := tt.expectedOutput.err; !errors.Is(err, expectedError) {
				t.Errorf("SearchBooks() got error %+v, expected %+v",
					err, expectedError)
			}
			if expectedBooks := tt.expectedOutput.books; !reflect.DeepEqual(books, expectedBooks) {
				t.Errorf("SearchBooks() got books %+v, expected %+v",
					books, expectedBooks)
			}
		})
	}
}

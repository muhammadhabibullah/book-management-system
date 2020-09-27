package mysql

import (
	"context"
	"errors"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"

	"book-management-system/entities/models"
)

func TestNewBookRepository(t *testing.T) {
	db := &gorm.DB{}

	got := NewBookRepository(db)
	expected := &bookRepository{
		db: db,
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("NewBookRepository returns %+v\n expected %+v",
			got, expected)
	}
}

func TestBookRepository_GetAll(t *testing.T) {
	type input struct {
		ctx context.Context
	}
	type output struct {
		books models.Books
		err   error
	}
	type mockConfig struct {
		expected output
		mock     sqlmock.Sqlmock
	}

	queryRgx := regexp.QuoteMeta("SELECT * FROM `books` WHERE `books`.`deleted_at` IS NULL")
	errDatabase := errors.New("error")

	tests := []struct {
		name           string
		givenInput     input
		expectedOutput output
		configureMock  func(mockConfig)
	}{
		{
			name: "success get all books",
			givenInput: input{
				ctx: context.TODO(),
			},
			expectedOutput: output{
				books: models.Books{
					{
						Model: gorm.Model{
							ID: 1,
						},
						Name: "Book",
						ISBN: "1234",
					},
				},
				err: nil,
			},
			configureMock: func(conf mockConfig) {
				rows := sqlmock.NewRows([]string{"id", "name", "isbn"})
				for _, book := range conf.expected.books {
					rows.AddRow(book.ID, book.Name, book.ISBN)
				}

				conf.mock.ExpectQuery(queryRgx).
					WillReturnRows(rows)
			},
		},
		{
			name: "no books found",
			givenInput: input{
				ctx: context.TODO(),
			},
			expectedOutput: output{
				books: models.Books{},
				err:   nil,
			},
			configureMock: func(conf mockConfig) {
				conf.mock.ExpectQuery(queryRgx).
					WillReturnRows(&sqlmock.Rows{})
			},
		},
		{
			name: "error database",
			expectedOutput: output{
				books: models.Books{},
				err:   errDatabase,
			},
			configureMock: func(conf mockConfig) {
				conf.mock.ExpectQuery(queryRgx).
					WillReturnError(conf.expected.err)
			},
		},
	}

	dbMock, mock, err := setupTestSuite()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer closeDB(dbMock)

	for _, tt := range tests {
		tt.configureMock(mockConfig{
			expected: tt.expectedOutput,
			mock:     mock,
		})

		repo := bookRepository{
			db: dbMock,
		}

		books, err := repo.GetAll(tt.givenInput.ctx)
		if expectedError := tt.expectedOutput.err; !errors.Is(err, expectedError) {
			t.Errorf("GetAll() got error: %v\nexpected: %v",
				err, expectedError)
		}
		if expectedBooks := tt.expectedOutput.books; err == nil && !reflect.DeepEqual(books, expectedBooks) {
			t.Errorf("GetAll() got books: %+v \nexpected: %+v",
				books, expectedBooks)
		}
	}
}

func TestBookRepository_CreateBook(t *testing.T) {
	type input struct {
		ctx  context.Context
		book *models.Book
	}
	type output struct {
		err error
	}
	type mockConfig struct {
		given    input
		expected output
		mock     sqlmock.Sqlmock
	}

	queryRgx := regexp.QuoteMeta("INSERT INTO `books` (`created_at`,`updated_at`,`deleted_at`,`name`,`isbn`) VALUES (?,?,?,?,?)")
	errDatabase := errors.New("error database")

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
					Name: "Book",
					ISBN: "1234",
				},
			},
			expectedOutput: output{
				err: nil,
			},
			configureMock: func(conf mockConfig) {
				conf.mock.ExpectBegin()
				conf.mock.ExpectExec(queryRgx).
					WithArgs(
						AnyTime{}, AnyTime{}, nil,
						conf.given.book.Name,
						conf.given.book.ISBN,
					).WillReturnResult(sqlmock.NewResult(1, 1))
				conf.mock.ExpectCommit()
			},
		},
		{
			name: "error create book",
			givenInput: input{
				ctx: context.TODO(),
				book: &models.Book{
					Name: "Book",
					ISBN: "1234",
				},
			},
			expectedOutput: output{
				err: errDatabase,
			},
			configureMock: func(conf mockConfig) {
				conf.mock.ExpectBegin()
				conf.mock.ExpectExec(queryRgx).
					WithArgs(
						AnyTime{}, AnyTime{}, nil,
						conf.given.book.Name,
						conf.given.book.ISBN,
					).WillReturnError(conf.expected.err)
				conf.mock.ExpectRollback()
			},
		},
	}

	dbMock, mock, err := setupTestSuite()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer closeDB(dbMock)

	for _, tt := range tests {
		tt.configureMock(mockConfig{
			given:    tt.givenInput,
			expected: tt.expectedOutput,
			mock:     mock,
		})

		repo := bookRepository{
			db: dbMock,
		}

		err := repo.CreateBook(tt.givenInput.ctx, tt.givenInput.book)
		if expectedError := tt.expectedOutput.err; !errors.Is(err, expectedError) {
			t.Errorf("CreateBook() got error: %v\nexpected: %v",
				err, expectedError)
		}
	}
}

func TestBookRepository_UpdateBook(t *testing.T) {
	type input struct {
		ctx  context.Context
		book *models.Book
	}
	type output struct {
		err error
	}
	type mockConfig struct {
		given    input
		expected output
		mock     sqlmock.Sqlmock
	}

	queryRgx := regexp.QuoteMeta("UPDATE `books` SET `updated_at`=?,`name`=?,`isbn`=? WHERE `id` = ?")
	errDatabase := errors.New("error")

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
					Model: gorm.Model{
						ID: 1,
					},
					Name: "Book",
					ISBN: "1234",
				},
			},
			expectedOutput: output{
				err: nil,
			},
			configureMock: func(conf mockConfig) {
				conf.mock.ExpectBegin()
				conf.mock.ExpectExec(queryRgx).
					WithArgs(
						AnyTime{},
						conf.given.book.Name,
						conf.given.book.ISBN,
						conf.given.book.ID,
					).WillReturnResult(sqlmock.NewResult(1, 1))
				conf.mock.ExpectCommit()
			},
		},
		{
			name: "error update book",
			givenInput: input{
				ctx: context.TODO(),
				book: &models.Book{
					Model: gorm.Model{
						ID: 1,
					},
					Name: "Book",
					ISBN: "1234",
				},
			},
			expectedOutput: output{
				err: errDatabase,
			},
			configureMock: func(conf mockConfig) {
				conf.mock.ExpectBegin()
				conf.mock.ExpectExec(queryRgx).
					WithArgs(
						AnyTime{},
						conf.given.book.Name,
						conf.given.book.ISBN,
						conf.given.book.ID,
					).WillReturnError(conf.expected.err)
				conf.mock.ExpectRollback()
			},
		},
	}

	dbMock, mock, err := setupTestSuite()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer closeDB(dbMock)

	for _, tt := range tests {
		tt.configureMock(mockConfig{
			given:    tt.givenInput,
			expected: tt.expectedOutput,
			mock:     mock,
		})

		repo := bookRepository{
			db: dbMock,
		}

		err := repo.UpdateBook(tt.givenInput.ctx, tt.givenInput.book)
		if expectedError := tt.expectedOutput.err; !errors.Is(err, expectedError) {
			t.Errorf("UpdateBook() got error: %v\nexpected: %v",
				err, expectedError)
		}
	}
}

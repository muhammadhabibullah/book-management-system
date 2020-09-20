package mysql

import (
	"errors"
	"reflect"
	"regexp"
	"testing"

	"book-management-system/entities/models"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"
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
	type output struct {
		books models.Books
		err   error
	}
	type confMock struct {
		output output
		mock   sqlmock.Sqlmock
	}

	queryRgx := regexp.QuoteMeta("SELECT * FROM `books` WHERE `books`.`deleted_at` IS NULL")
	errDatabase := errors.New("error")

	tests := []struct {
		name          string
		output        output
		configureMock func(confMock)
	}{
		{
			name: "success get all books",
			output: output{
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
			configureMock: func(conf confMock) {
				rows := sqlmock.NewRows([]string{"id", "name", "isbn"}).
					AddRow(1, "Book", "1234")
				conf.mock.ExpectQuery(queryRgx).
					WillReturnRows(rows)
			},
		},
		{
			name: "no books found",
			output: output{
				books: models.Books{},
				err:   nil,
			},
			configureMock: func(conf confMock) {
				conf.mock.ExpectQuery(queryRgx).
					WillReturnRows(&sqlmock.Rows{})
			},
		},
		{
			name: "error database",
			output: output{
				books: models.Books{},
				err:   errDatabase,
			},
			configureMock: func(conf confMock) {
				conf.mock.ExpectQuery(queryRgx).
					WillReturnError(errDatabase)
			},
		},
	}

	dbMock, mock, err := setupTestSuite()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	for _, tt := range tests {
		tt.configureMock(confMock{
			output: tt.output,
			mock:   mock,
		})

		repo := bookRepository{
			db: dbMock,
		}

		books, err := repo.GetAll()
		if expectedErr := tt.output.err; !errors.Is(err, tt.output.err) {
			t.Errorf("GetAll() got error: %v\nexpected: %v",
				err, expectedErr)
		}
		if expectedBooks := tt.output.books; err == nil && !reflect.DeepEqual(books, expectedBooks) {
			t.Errorf("GetAll() got books: %+v with type %T\nexpected: %+v with type %T",
				books, books, expectedBooks, expectedBooks)
		}
	}
}

func TestBookRepository_CreateBook(t *testing.T) {
	type input struct {
		book *models.Book
	}
	type output struct {
		err error
	}
	type confMock struct {
		input  input
		output output
		mock   sqlmock.Sqlmock
	}

	queryRgx := regexp.QuoteMeta("INSERT INTO `books` (`created_at`,`updated_at`,`deleted_at`,`name`,`isbn`) VALUES (?,?,?,?,?)")
	errDatabase := errors.New("error")

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
					Name: "Book",
					ISBN: "1234",
				},
			},
			output: output{
				err: nil,
			},
			configureMock: func(conf confMock) {
				conf.mock.ExpectBegin()
				conf.mock.ExpectExec(queryRgx).
					WithArgs(
						AnyTime{}, AnyTime{}, nil,
						conf.input.book.Name,
						conf.input.book.ISBN,
					).WillReturnResult(sqlmock.NewResult(1, 1))
				conf.mock.ExpectCommit()
			},
		},
		{
			name: "error create book",
			input: input{
				book: &models.Book{
					Name: "Book",
					ISBN: "1234",
				},
			},
			output: output{
				err: errDatabase,
			},
			configureMock: func(conf confMock) {
				conf.mock.ExpectBegin()
				conf.mock.ExpectExec(queryRgx).
					WithArgs(
						AnyTime{}, AnyTime{}, nil,
						conf.input.book.Name,
						conf.input.book.ISBN,
					).WillReturnError(errDatabase)
				conf.mock.ExpectRollback()
			},
		},
	}

	dbMock, mock, err := setupTestSuite()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	for _, tt := range tests {
		tt.configureMock(confMock{
			input:  tt.input,
			output: tt.output,
			mock:   mock,
		})

		repo := bookRepository{
			db: dbMock,
		}

		err := repo.CreateBook(tt.input.book)
		if expectedErr := tt.output.err; !errors.Is(err, tt.output.err) {
			t.Errorf("CreateBook() got error: %v\nexpected: %v",
				err, expectedErr)
		}
	}
}

func TestBookRepository_UpdateBook(t *testing.T) {
	type input struct {
		book *models.Book
	}
	type output struct {
		err error
	}
	type confMock struct {
		input  input
		output output
		mock   sqlmock.Sqlmock
	}

	queryRgx := regexp.QuoteMeta("UPDATE `books` SET `updated_at`=?,`name`=?,`isbn`=? WHERE `id` = ?")
	errDatabase := errors.New("error")

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
					Model: gorm.Model{
						ID: 1,
					},
					Name: "Book",
					ISBN: "1234",
				},
			},
			output: output{
				err: nil,
			},
			configureMock: func(conf confMock) {
				conf.mock.ExpectBegin()
				conf.mock.ExpectExec(queryRgx).
					WithArgs(
						AnyTime{},
						conf.input.book.Name,
						conf.input.book.ISBN,
						conf.input.book.ID,
					).WillReturnResult(sqlmock.NewResult(1, 1))
				conf.mock.ExpectCommit()
			},
		},
		{
			name: "error update book",
			input: input{
				book: &models.Book{
					Model: gorm.Model{
						ID: 1,
					},
					Name: "Book",
					ISBN: "1234",
				},
			},
			output: output{
				err: errDatabase,
			},
			configureMock: func(conf confMock) {
				conf.mock.ExpectBegin()
				conf.mock.ExpectExec(queryRgx).
					WithArgs(
						AnyTime{},
						conf.input.book.Name,
						conf.input.book.ISBN,
						conf.input.book.ID,
					).WillReturnError(errDatabase)
				conf.mock.ExpectRollback()
			},
		},
	}

	dbMock, mock, err := setupTestSuite()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	for _, tt := range tests {
		tt.configureMock(confMock{
			input:  tt.input,
			output: tt.output,
			mock:   mock,
		})

		repo := bookRepository{
			db: dbMock,
		}

		err := repo.UpdateBook(tt.input.book)
		if expectedErr := tt.output.err; !errors.Is(err, tt.output.err) {
			t.Errorf("UpdateBook() got error: %v\nexpected: %v",
				err, expectedErr)
		}
	}
}

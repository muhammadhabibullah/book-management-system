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
}

func TestBookService_CreateBook(t *testing.T) {
	type input struct {
		book *models.Book
	}
	type output struct {
		err error
	}
	type confMock struct {
		input input
		mock  mocks.MockBookRepository
	}

	tests := []struct {
		name          string
		input         input
		output        output
		configureMock func(confMock)
	}{
		{
			name: "create book",
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
				conf.mock.EXPECT().
					CreateBook(conf.input.book).
					Return(nil)
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
				input: tt.input,
				mock:  *bookRepoMock,
			})

			err := bookService.CreateBook(tt.input.book)
			if !errors.Is(err, tt.output.err) {
				t.Errorf("unexpected error %+v, expected %+v",
					err, tt.output.err)
			}
		})
	}
}

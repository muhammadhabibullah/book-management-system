package services

import (
	"errors"
	"reflect"
	"testing"

	"book-management-system/repositories"
)

func TestInitServices(t *testing.T) {
	repo := &repositories.Repository{}

	got := Init(repo)
	expected := &Services{
		BookService:   NewBookService(repo),
		MemberService: NewMemberService(repo),
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Init returns %+v\n expected %+v",
			got, expected)
	}
}

var errRepository = errors.New("repository error")

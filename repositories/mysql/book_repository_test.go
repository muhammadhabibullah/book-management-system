package mysql

import (
	"reflect"
	"testing"

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

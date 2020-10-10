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

func TestNewMemberRepository(t *testing.T) {
	db := &gorm.DB{}

	got := NewMemberRepository(db)
	expected := &memberRepository{
		db: db,
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("NewMemberRepository returns %+v\n expected %+v",
			got, expected)
	}
}

func TestMemberRepositoryGetAll(t *testing.T) {
	type input struct {
		ctx context.Context
	}
	type output struct {
		members models.Members
		err     error
	}
	type mockConfig struct {
		expected output
		mock     sqlmock.Sqlmock
	}

	queryRgx := regexp.QuoteMeta("SELECT * FROM `members` WHERE `members`.`deleted_at` IS NULL")

	tests := []struct {
		name           string
		givenInput     input
		expectedOutput output
		configureMock  func(mockConfig)
	}{
		{
			name: "success get all members",
			givenInput: input{
				ctx: context.TODO(),
			},
			expectedOutput: output{
				members: models.Members{
					{
						Model: gorm.Model{
							ID: 1,
						},
						Name: "",
					},
				},
				err: nil,
			},
			configureMock: func(conf mockConfig) {
				rows := sqlmock.NewRows([]string{"id", "name"})
				for _, member := range conf.expected.members {
					rows.AddRow(member.ID, member.Name)
				}

				conf.mock.ExpectQuery(queryRgx).
					WillReturnRows(rows)
			},
		},
		{
			name: "no members found",
			givenInput: input{
				ctx: context.TODO(),
			},
			expectedOutput: output{
				members: models.Members{},
				err:     nil,
			},
			configureMock: func(conf mockConfig) {
				conf.mock.ExpectQuery(queryRgx).
					WillReturnRows(&sqlmock.Rows{})
			},
		},
		{
			name: "error database",
			expectedOutput: output{
				members: models.Members{},
				err:     errDatabase,
			},
			configureMock: func(conf mockConfig) {
				conf.mock.ExpectQuery(queryRgx).
					WillReturnError(conf.expected.err)
			},
		},
	}

	dbMock, mock, err := setupTestSuite()
	if err != nil {
		t.Fatal(err)
	}
	defer closeDB(dbMock)

	for _, tt := range tests {
		tt.configureMock(mockConfig{
			expected: tt.expectedOutput,
			mock:     mock,
		})

		repo := memberRepository{
			db: dbMock,
		}

		members, err := repo.GetAll(tt.givenInput.ctx)
		if expectedError := tt.expectedOutput.err; !errors.Is(err, expectedError) {
			t.Errorf("GetAll() got error: %v\nexpected: %v",
				err, expectedError)
		}
		if expectedMembers := tt.expectedOutput.members; err == nil && !reflect.DeepEqual(members, expectedMembers) {
			t.Errorf("GetAll() got members: %+v \nexpected: %+v",
				members, expectedMembers)
		}
	}
}

func TestMemberRepositoryCreateMember(t *testing.T) {
	type input struct {
		ctx    context.Context
		member *models.Member
	}
	type output struct {
		err error
	}
	type mockConfig struct {
		given    input
		expected output
		mock     sqlmock.Sqlmock
	}

	queryRgx := regexp.QuoteMeta("INSERT INTO `members` (`created_at`,`updated_at`,`deleted_at`,`name`) VALUES (?,?,?,?)")

	tests := []struct {
		name           string
		givenInput     input
		expectedOutput output
		configureMock  func(mockConfig)
	}{
		{
			name: "success create member",
			givenInput: input{
				ctx: context.TODO(),
				member: &models.Member{
					Name: "",
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
						conf.given.member.Name,
					).WillReturnResult(sqlmock.NewResult(1, 1))
				conf.mock.ExpectCommit()
			},
		},
		{
			name: "error create member",
			givenInput: input{
				ctx: context.TODO(),
				member: &models.Member{
					Name: "",
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
						conf.given.member.Name,
					).WillReturnError(conf.expected.err)
				conf.mock.ExpectRollback()
			},
		},
	}

	dbMock, mock, err := setupTestSuite()
	if err != nil {
		t.Fatal(err)
	}
	defer closeDB(dbMock)

	for _, tt := range tests {
		tt.configureMock(mockConfig{
			given:    tt.givenInput,
			expected: tt.expectedOutput,
			mock:     mock,
		})

		repo := memberRepository{
			db: dbMock,
		}

		err := repo.CreateMember(tt.givenInput.ctx, tt.givenInput.member)
		if expectedError := tt.expectedOutput.err; !errors.Is(err, expectedError) {
			t.Errorf("CreateMember() got error: %v\nexpected: %v",
				err, expectedError)
		}
	}
}

func TestMemberRepositoryUpdateMember(t *testing.T) {
	type input struct {
		ctx    context.Context
		member *models.Member
	}
	type output struct {
		err error
	}
	type mockConfig struct {
		given    input
		expected output
		mock     sqlmock.Sqlmock
	}

	queryRgx := regexp.QuoteMeta("UPDATE `members` SET `updated_at`=?,`name`=? WHERE `id` = ?")

	tests := []struct {
		name           string
		givenInput     input
		expectedOutput output
		configureMock  func(mockConfig)
	}{
		{
			name: "success update member",
			givenInput: input{
				ctx: context.TODO(),
				member: &models.Member{
					Model: gorm.Model{
						ID: 1,
					},
					Name: "Updated Member",
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
						conf.given.member.Name,
						conf.given.member.ID,
					).WillReturnResult(sqlmock.NewResult(1, 1))
				conf.mock.ExpectCommit()
			},
		},
		{
			name: "error update member",
			givenInput: input{
				ctx: context.TODO(),
				member: &models.Member{
					Model: gorm.Model{
						ID: 1,
					},
					Name: "Updated Member",
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
						conf.given.member.Name,
						conf.given.member.ID,
					).WillReturnError(conf.expected.err)
				conf.mock.ExpectRollback()
			},
		},
	}

	dbMock, mock, err := setupTestSuite()
	if err != nil {
		t.Fatal(err)
	}
	defer closeDB(dbMock)

	for _, tt := range tests {
		tt.configureMock(mockConfig{
			given:    tt.givenInput,
			expected: tt.expectedOutput,
			mock:     mock,
		})

		repo := memberRepository{
			db: dbMock,
		}

		err := repo.UpdateMember(tt.givenInput.ctx, tt.givenInput.member)
		if expectedError := tt.expectedOutput.err; !errors.Is(err, expectedError) {
			t.Errorf("UpdateMember() got error: %v\nexpected: %v",
				err, expectedError)
		}
	}
}

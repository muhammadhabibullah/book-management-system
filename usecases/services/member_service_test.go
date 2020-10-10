package services

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"book-management-system/entities/models"
	mySqlMocks "book-management-system/mocks/repositories/mysql"
	"book-management-system/repositories"
	"book-management-system/repositories/mysql"
)

func TestNewMemberService(t *testing.T) {
	mySQLMemberRepo := mysql.NewMemberRepository(nil)
	repo := &repositories.Repository{
		MySQLMemberRepository: mySQLMemberRepo,
	}

	got := NewMemberService(repo)
	expected := &memberService{
		MySQLMemberRepository: mySQLMemberRepo,
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("NewMemberService returns %+v\n expected %+v",
			got, expected)
	}
	if _, ok := got.(MemberService); !ok {
		t.Errorf("NewMemberService returns object not implements MemberService")
	}
}

func TestMemberServiceCreateMember(t *testing.T) {
	type input struct {
		ctx    context.Context
		member *models.Member
	}
	type output struct {
		err error
	}
	type mockConfig struct {
		given               input
		expected            output
		mySQLMemberRepoMock *mySqlMocks.MockMemberRepository
	}

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
					Name: "John Lennon",
				},
			},
			expectedOutput: output{
				err: nil,
			},
			configureMock: func(conf mockConfig) {
				conf.mySQLMemberRepoMock.EXPECT().
					CreateMember(gomock.Any(), conf.given.member).
					Return(conf.expected.err)
			},
		},
		{
			name: "failed create member",
			givenInput: input{
				ctx: context.TODO(),
				member: &models.Member{
					Name: "John Lennon",
				},
			},
			expectedOutput: output{
				err: errRepository,
			},
			configureMock: func(conf mockConfig) {
				conf.mySQLMemberRepoMock.EXPECT().
					CreateMember(gomock.Any(), conf.given.member).
					Return(conf.expected.err)
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mySQLMemberRepoMock := mySqlMocks.NewMockMemberRepository(ctrl)

			memberService := &memberService{
				MySQLMemberRepository: mySQLMemberRepoMock,
			}

			tt.configureMock(mockConfig{
				given:               tt.givenInput,
				expected:            tt.expectedOutput,
				mySQLMemberRepoMock: mySQLMemberRepoMock,
			})

			err := memberService.CreateMember(tt.givenInput.ctx, tt.givenInput.member)
			if expectedError := tt.expectedOutput.err; !errors.Is(err, expectedError) {
				t.Errorf("CreateMember() got error %+v, expected %+v",
					err, expectedError)
			}
		})
	}
}

func TestMemberServiceGetMembers(t *testing.T) {
	type input struct {
		ctx context.Context
	}
	type output struct {
		members models.Members
		err     error
	}
	type mockConfig struct {
		given               input
		expected            output
		mySQLMemberRepoMock *mySqlMocks.MockMemberRepository
	}

	tests := []struct {
		name           string
		givenInput     input
		expectedOutput output
		configureMock  func(mockConfig)
	}{
		{
			name: "success get members",
			givenInput: input{
				ctx: context.TODO(),
			},
			expectedOutput: output{
				members: models.Members{
					{
						Name: "John Lennon",
					},
				},
				err: nil,
			},
			configureMock: func(conf mockConfig) {
				conf.mySQLMemberRepoMock.EXPECT().
					GetAll(gomock.Any()).
					Return(
						conf.expected.members,
						conf.expected.err,
					)
			},
		},
		{
			name: "failed get members",
			givenInput: input{
				ctx: context.TODO(),
			},
			expectedOutput: output{
				members: models.Members{},
				err:     errRepository,
			},
			configureMock: func(conf mockConfig) {
				conf.mySQLMemberRepoMock.EXPECT().
					GetAll(gomock.Any()).
					Return(
						conf.expected.members,
						conf.expected.err,
					)
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mySQLMemberRepoMock := mySqlMocks.NewMockMemberRepository(ctrl)

			memberService := &memberService{
				MySQLMemberRepository: mySQLMemberRepoMock,
			}

			tt.configureMock(mockConfig{
				given:               tt.givenInput,
				expected:            tt.expectedOutput,
				mySQLMemberRepoMock: mySQLMemberRepoMock,
			})

			members, err := memberService.GetMembers(tt.givenInput.ctx)
			if expectedError := tt.expectedOutput.err; !errors.Is(err, expectedError) {
				t.Errorf("GetMembers() got error %+v, expected %+v",
					err, expectedError)
			}
			if expectedMembers := tt.expectedOutput.members; !reflect.DeepEqual(members, expectedMembers) {
				t.Errorf("GetMembers() got members %+v, expected %+v",
					members, expectedMembers)
			}
		})
	}
}

func TestMemberServiceUpdateMember(t *testing.T) {
	type input struct {
		ctx    context.Context
		member *models.Member
	}
	type output struct {
		err error
	}
	type mockConfig struct {
		given               input
		expected            output
		mySQLMemberRepoMock *mySqlMocks.MockMemberRepository
	}

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
					Name: "John Lennon",
				},
			},
			expectedOutput: output{
				err: nil,
			},
			configureMock: func(conf mockConfig) {
				conf.mySQLMemberRepoMock.EXPECT().
					UpdateMember(gomock.Any(), conf.given.member).
					Return(conf.expected.err)
			},
		},
		{
			name: "failed update member",
			givenInput: input{
				ctx: context.TODO(),
				member: &models.Member{
					Name: "John Lennon",
				},
			},
			expectedOutput: output{
				err: errRepository,
			},
			configureMock: func(conf mockConfig) {
				conf.mySQLMemberRepoMock.EXPECT().
					UpdateMember(gomock.Any(), conf.given.member).
					Return(conf.expected.err)
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mySQLMemberRepoMock := mySqlMocks.NewMockMemberRepository(ctrl)

			memberService := &memberService{
				MySQLMemberRepository: mySQLMemberRepoMock,
			}

			tt.configureMock(mockConfig{
				given:               tt.givenInput,
				expected:            tt.expectedOutput,
				mySQLMemberRepoMock: mySQLMemberRepoMock,
			})

			err := memberService.UpdateMember(tt.givenInput.ctx, tt.givenInput.member)
			if expectedError := tt.expectedOutput.err; !errors.Is(err, expectedError) {
				t.Errorf("UpdateMember() got error %+v, expected %+v",
					err, expectedError)
			}
		})
	}
}

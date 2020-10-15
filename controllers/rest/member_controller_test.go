package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"

	"book-management-system/controllers/rest/middlewares"
	"book-management-system/controllers/rest/responses"
	"book-management-system/entities/models"
	mocks "book-management-system/mocks/services"
	"book-management-system/repositories"
	"book-management-system/usecases"
	"book-management-system/usecases/services"
)

func TestNewMemberController(t *testing.T) {
	repo := &repositories.Repository{}
	memberService := services.NewMemberService(repo)
	usecase := &usecases.UseCase{
		Service: &services.Services{
			MemberService: memberService,
		},
	}

	route := mux.NewRouter()
	m := &middlewares.Middleware{}
	got := NewMemberController(route, usecase, m)
	expected := &MemberController{
		memberService: memberService,
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("NewMemberController returns %+v\n expected %+v",
			got, expected)
	}
}

const (
	v1MemberURL = "/v1/member"
)

func TestMemberControllerCreateMember(t *testing.T) {
	type input struct {
		valid              bool
		ctx                context.Context
		requestBody        *models.Member
		invalidRequestBody models.Members
	}
	type output struct {
		responseBody interface{}
	}
	type mockConfig struct {
		given input
		mock  *mocks.MockMemberService
	}

	tests := []struct {
		name           string
		givenInput     input
		expectedOutput output
		configureMock  func(mockConfig)
	}{
		{
			name: "failed: invalid request body",
			givenInput: input{
				valid: false,
				ctx:   context.TODO(),
				invalidRequestBody: models.Members{
					{
						Name: "",
					},
				},
			},
			expectedOutput: output{
				responseBody: responses.ErrorResponse{
					"error": "Invalid request payload",
				},
			},
			configureMock: func(mockConfig) {
				// do nothing
			},
		},
		{
			name: "failed: create member service returns error",
			givenInput: input{
				valid: true,
				ctx:   context.TODO(),
				requestBody: &models.Member{
					Name: "",
				},
			},
			expectedOutput: output{
				responseBody: responses.ErrorResponse{
					"error": fmt.Sprintf("Failed create member: %s", errService.Error()),
				},
			},
			configureMock: func(conf mockConfig) {
				conf.mock.EXPECT().
					CreateMember(conf.given.ctx, conf.given.requestBody).
					Return(errService)
			},
		},
		{
			name: "success: create member",
			givenInput: input{
				valid: true,
				ctx:   context.TODO(),
				requestBody: &models.Member{
					Name: "",
				},
			},
			expectedOutput: output{
				responseBody: models.Member{
					Name: "",
				},
			},
			configureMock: func(conf mockConfig) {
				conf.mock.EXPECT().
					CreateMember(conf.given.ctx, conf.given.requestBody).
					Return(nil)
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var marshalledRequestBody []byte
			if tt.givenInput.valid {
				marshalledRequestBody, _ = json.Marshal(tt.givenInput.requestBody)
			} else {
				marshalledRequestBody, _ = json.Marshal(tt.givenInput.invalidRequestBody)
			}

			req, _ := http.NewRequest(
				http.MethodPost,
				v1MemberURL,
				bytes.NewBuffer(marshalledRequestBody),
			)
			resp := httptest.NewRecorder()

			memberServiceMock := mocks.NewMockMemberService(ctrl)
			tt.configureMock(mockConfig{
				given: tt.givenInput,
				mock:  memberServiceMock,
			})

			memberController := &MemberController{
				memberService: memberServiceMock,
			}

			memberController.CreateMember(resp, req)

			got := resp.Body.String()
			expected, _ := json.Marshal(tt.expectedOutput.responseBody)
			if got != string(expected) {
				t.Errorf("CreateMember() got response body %s\n expected %s",
					got, string(expected))
			}
		})
	}
}

func TestMemberControllerGetMember(t *testing.T) {
	type input struct {
		ctx            context.Context
		httpRequestURL string
	}
	type output struct {
		responseBody interface{}
	}
	type mockConfig struct {
		given    input
		expected output
		mock     *mocks.MockMemberService
	}

	tests := []struct {
		name           string
		givenInput     input
		expectedOutput output
		configureMock  func(mockConfig)
	}{
		{
			name: "failed: get members service returns error",
			givenInput: input{
				ctx:            context.TODO(),
				httpRequestURL: v1MemberURL,
			},
			expectedOutput: output{
				responseBody: responses.ErrorResponse{
					"error": fmt.Sprintf("Failed get members: %s", errService.Error()),
				},
			},
			configureMock: func(conf mockConfig) {
				conf.mock.EXPECT().
					GetMembers(conf.given.ctx).
					Return(models.Members{}, errService)
			},
		},
		{
			name: "success: get members",
			givenInput: input{
				ctx:            context.TODO(),
				httpRequestURL: v1MemberURL,
			},
			expectedOutput: output{
				responseBody: models.Members{
					{
						Name: "",
					},
				},
			},
			configureMock: func(conf mockConfig) {
				conf.mock.EXPECT().
					GetMembers(conf.given.ctx).
					Return(conf.expected.responseBody, nil)
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(
				http.MethodGet,
				tt.givenInput.httpRequestURL,
				nil,
			)
			resp := httptest.NewRecorder()

			memberServiceMock := mocks.NewMockMemberService(ctrl)
			tt.configureMock(mockConfig{
				given:    tt.givenInput,
				expected: tt.expectedOutput,
				mock:     memberServiceMock,
			})

			memberController := &MemberController{
				memberService: memberServiceMock,
			}

			memberController.GetMembers(resp, req)

			got := resp.Body.String()
			expected, _ := json.Marshal(tt.expectedOutput.responseBody)
			if got != string(expected) {
				t.Errorf("GetMembers() got response body %s\n expected %s",
					got, string(expected))
			}
		})
	}
}

func TestMemberControllerUpdateMember(t *testing.T) {
	type input struct {
		valid              bool
		ctx                context.Context
		requestBody        *models.Member
		invalidRequestBody models.Members
	}
	type output struct {
		responseBody interface{}
	}
	type confMock struct {
		given input
		mock  *mocks.MockMemberService
	}

	tests := []struct {
		name           string
		givenInput     input
		expectedOutput output
		configureMock  func(confMock)
	}{
		{
			name: "failed: invalid request body",
			givenInput: input{
				valid: false,
				ctx:   context.TODO(),
				invalidRequestBody: models.Members{
					{
						Name: "",
					},
				},
			},
			expectedOutput: output{
				responseBody: responses.ErrorResponse{
					"error": "Invalid request payload",
				},
			},
			configureMock: func(confMock) {
				// do nothing
			},
		},
		{
			name: "failed: update member service returns error",
			givenInput: input{
				valid: true,
				ctx:   context.TODO(),
				requestBody: &models.Member{
					Name: "",
				},
			},
			expectedOutput: output{
				responseBody: responses.ErrorResponse{
					"error": fmt.Sprintf("Failed update member: %s", errService.Error()),
				},
			},
			configureMock: func(conf confMock) {
				conf.mock.EXPECT().
					UpdateMember(conf.given.ctx, conf.given.requestBody).
					Return(errService)
			},
		},
		{
			name: "success: update member",
			givenInput: input{
				valid: true,
				ctx:   context.TODO(),
				requestBody: &models.Member{
					Name: "",
				},
			},
			expectedOutput: output{
				responseBody: models.Member{
					Name: "",
				},
			},
			configureMock: func(conf confMock) {
				conf.mock.EXPECT().
					UpdateMember(conf.given.ctx, conf.given.requestBody).
					Return(nil)
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var marshalledRequestBody []byte
			if tt.givenInput.valid {
				marshalledRequestBody, _ = json.Marshal(tt.givenInput.requestBody)
			} else {
				marshalledRequestBody, _ = json.Marshal(tt.givenInput.invalidRequestBody)
			}

			req, _ := http.NewRequest(
				http.MethodPut,
				v1MemberURL,
				bytes.NewBuffer(marshalledRequestBody),
			)
			resp := httptest.NewRecorder()

			memberServiceMock := mocks.NewMockMemberService(ctrl)
			tt.configureMock(confMock{
				given: tt.givenInput,
				mock:  memberServiceMock,
			})

			memberController := &MemberController{
				memberService: memberServiceMock,
			}

			memberController.UpdateMember(resp, req)

			got := resp.Body.String()
			expected, _ := json.Marshal(tt.expectedOutput.responseBody)
			if got != string(expected) {
				t.Errorf("UpdateMember() got response body %s\n expected %s",
					got, string(expected))
			}
		})
	}
}

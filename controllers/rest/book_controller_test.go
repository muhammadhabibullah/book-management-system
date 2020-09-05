package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"

	"book-management-system/entities/models"
	mocks "book-management-system/mocks/services"
)

func TestNewBookController(t *testing.T) {}

func TestBookController_CreateBook(t *testing.T) {
	type input struct {
		requestBody *models.Book
	}
	type output struct {
		responseBody interface{}
	}
	type confMock struct {
		input input
		mock  *mocks.MockBookService
	}

	tests := []struct {
		name          string
		input         input
		output        output
		configureMock func(confMock)
	}{
		// TODO
		//{
		//	name: "failed: invalid request body",
		//	input: input{
		//		requestBody: nil,
		//	},
		//	output: output{
		//		responseBody: responses.ErrorResponse{
		//			"error": "Invalid request payload",
		//		},
		//	},
		//	configureMock: func(confMock) {},
		//},
		{
			name: "success: create book",
			input: input{
				requestBody: &models.Book{
					Name: "C++",
					ISBN: "1234",
				},
			},
			output: output{
				responseBody: models.Book{
					Name: "C++",
					ISBN: "1234",
				},
			},
			configureMock: func(conf confMock) {
				conf.mock.EXPECT().
					CreateBook(gomock.Any()).
					SetArg(0, *conf.input.requestBody).
					Return(nil)
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			marshalledRequestBody, _ := json.Marshal(tt.input.requestBody)
			req, _ := http.NewRequest(
				"POST",
				"/v1/book",
				bytes.NewBuffer(marshalledRequestBody),
			)
			resp := httptest.NewRecorder()

			bookServiceMock := mocks.NewMockBookService(ctrl)
			tt.configureMock(confMock{
				input: tt.input,
				mock:  bookServiceMock,
			})

			bookController := &BookController{
				bookService: bookServiceMock,
			}

			bookController.CreateBook(resp, req)

			got := resp.Body.String()
			expected, _ := json.Marshal(tt.output.responseBody)
			if got != string(expected) {
				t.Errorf("got response body %s\n expected %s",
					got, string(expected))
			}
		})
	}
}

package httpserver

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cronnoss/bookshop-home-task/internal/app/domain"
	mock_service "github.com/cronnoss/bookshop-home-task/internal/app/services/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestHTTPServer_SignUp_Success(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUserRepository, user domain.User)

	testTable := []struct {
		name         string
		inputBody    string
		inputUser    domain.User
		mockBehavior mockBehavior
		expectedCode int
	}{
		{
			name:      "OK",
			inputBody: `{"username" : "test", "password" : "qwerty"}`,
			inputUser: domain.User{
				Username: "test",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockUserRepository, user domain.User) {
				s.EXPECT().CreateUser(gomock.Any(), gomock.AssignableToTypeOf(domain.User{})).Return(domain.User{ID: 1}, nil)
			},
			expectedCode: 200,
		},
		{
			name:         "Empty fields",
			inputBody:    `{"username" : "test"}`,
			mockBehavior: func(s *mock_service.MockUserRepository, user domain.User) {},
			expectedCode: 400,
		},
		{
			name:      "Service failure",
			inputBody: `{"username" : "test", "password" : "qwerty"}`,
			inputUser: domain.User{
				Username: "test",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockUserRepository, user domain.User) {
				s.EXPECT().CreateUser(gomock.Any(), gomock.AssignableToTypeOf(domain.User{})).Return(domain.User{},
					errors.New("internal error"))
			},
			expectedCode: 500,
		},
		{
			name:         "Empty body",
			inputBody:    `{}`,
			mockBehavior: func(s *mock_service.MockUserRepository, user domain.User) {},
			expectedCode: 400,
		},
		{
			name:         "Invalid body",
			inputBody:    `{"username" : "test" "":"" "password" : "qwerty"}`,
			mockBehavior: func(s *mock_service.MockUserRepository, user domain.User) {},
			expectedCode: 400,
		},
		{
			name:         "Invalid JSON",
			inputBody:    `{"username" : "test", "password" : "qwerty"`, // missing closing bracket
			mockBehavior: func(s *mock_service.MockUserRepository, user domain.User) {},
			expectedCode: 400,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockUserRepository(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			// Test server
			httpServer := NewHTTPServer(auth, nil, nil, nil, nil)
			req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBufferString(testCase.inputBody))

			// Test request
			w := httptest.NewRecorder()

			// Perform request
			httpServer.SignUp(w, req)

			// Assert
			require.Equal(t, testCase.expectedCode, w.Code)
		})
	}
}

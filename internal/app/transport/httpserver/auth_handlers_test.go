package httpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cronnoss/bookshop-home-task/internal/app/domain"
	"github.com/cronnoss/bookshop-home-task/internal/app/transport/httpserver/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignUp_Success(t *testing.T) {
	userServiceMock := mocks.NewUserService(t)

	userServiceMock.On("CreateUser", mock.Anything, mock.Anything).Return(domain.User{}, nil)

	httpServer := NewHTTPServer(userServiceMock, nil, nil, nil, nil)

	reqBody := AuthRequest{
		Username: "testuser",
		Password: "password123",
	}
	reqBodyJSON, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/signup",
		bytes.NewBuffer(reqBodyJSON))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	hashedPassword, _ := hashPassword(reqBody.Password)
	user, _ := toDomainUser(reqBody.Username, hashedPassword)

	userServiceMock.On("CreateUser", mock.Anything, mock.Anything).Return(user, nil)

	httpServer.SignUp(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	userServiceMock.AssertNumberOfCalls(t, "CreateUser", 1)
}

func TestSignUp_Validate(t *testing.T) {
	userServiceMock := mocks.NewUserService(t)

	userServiceMock.AssertExpectations(t)

	httpServer := NewHTTPServer(userServiceMock, nil, nil, nil, nil)

	t.Run("empty username", func(t *testing.T) {
		reqBody := AuthRequest{
			Username: "",
			Password: "password123",
		}
		reqBodyJSON, _ := json.Marshal(reqBody)

		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/signup",
			bytes.NewBuffer(reqBodyJSON))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		httpServer.SignUp(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		userServiceMock.AssertNumberOfCalls(t, "CreateUser", 0)
	})

	t.Run("empty password", func(t *testing.T) {
		reqBody := AuthRequest{
			Username: "testuser",
			Password: "",
		}
		reqBodyJSON, _ := json.Marshal(reqBody)

		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/signup",
			bytes.NewBuffer(reqBodyJSON))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		httpServer.SignUp(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		userServiceMock.AssertNumberOfCalls(t, "CreateUser", 0)
	})

	t.Run("invalid json", func(t *testing.T) {
		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/signup",
			bytes.NewBuffer([]byte("{invalid}")))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		httpServer.SignUp(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		userServiceMock.AssertNumberOfCalls(t, "CreateUser", 0)
	})
}

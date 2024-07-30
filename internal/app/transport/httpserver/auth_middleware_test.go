package httpserver

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cronnoss/bookshop-home-task/internal/app/domain"
	"github.com/cronnoss/bookshop-home-task/internal/app/transport/httpserver/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckAdmin_ValidAdminToken(t *testing.T) {
	userServiceMock := mocks.NewUserService(t)
	tokenServiceMock := mocks.NewTokenService(t)
	httpServer := NewHTTPServer(userServiceMock, tokenServiceMock, nil, nil, nil)

	req, err := http.NewRequestWithContext(context.Background(), "", "", nil)
	require.NoError(t, err)
	req.Header.Set(AuthorizationHeader, BearerPrefix+"valid-admin-token")

	user := domain.User{Username: "admin", Admin: true}
	tokenServiceMock.On("GetUser", "valid-admin-token").Return(user, nil)

	rr := httptest.NewRecorder()
	handler := httpServer.CheckAdmin(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestCheckAdmin_InvalidToken(t *testing.T) {
	userServiceMock := mocks.NewUserService(t)
	tokenServiceMock := mocks.NewTokenService(t)
	httpServer := NewHTTPServer(userServiceMock, tokenServiceMock, nil, nil, nil)

	req, err := http.NewRequestWithContext(context.Background(), "", "", nil)
	require.NoError(t, err)
	req.Header.Set(AuthorizationHeader, BearerPrefix+"invalid-token")

	tokenServiceMock.On("GetUser", "invalid-token").Return(domain.User{}, errors.New("invalid token"))

	rr := httptest.NewRecorder()
	handler := httpServer.CheckAdmin(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCheckAdmin_NotAdminUser(t *testing.T) {
	userServiceMock := mocks.NewUserService(t)
	tokenServiceMock := mocks.NewTokenService(t)
	httpServer := NewHTTPServer(userServiceMock, tokenServiceMock, nil, nil, nil)

	req, err := http.NewRequestWithContext(context.Background(), "", "", nil)
	require.NoError(t, err)
	req.Header.Set(AuthorizationHeader, BearerPrefix+"valid-user-token")

	user := domain.User{Username: "user", Admin: false}
	tokenServiceMock.On("GetUser", "valid-user-token").Return(user, nil)

	rr := httptest.NewRecorder()
	handler := httpServer.CheckAdmin(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestCheckAuthorizedUser_ValidToken(t *testing.T) {
	userServiceMock := mocks.NewUserService(t)
	tokenServiceMock := mocks.NewTokenService(t)
	httpServer := NewHTTPServer(userServiceMock, tokenServiceMock, nil, nil, nil)

	req, err := http.NewRequestWithContext(context.Background(), "", "", nil)
	require.NoError(t, err)
	req.Header.Set(AuthorizationHeader, BearerPrefix+"valid-user-token")

	user := domain.User{Username: "user", Admin: false}
	tokenServiceMock.On("GetUser", "valid-user-token").Return(user, nil)

	rr := httptest.NewRecorder()
	handler := httpServer.CheckAuthorizedUser(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestCheckAuthorizedUser_InvalidToken(t *testing.T) {
	userServiceMock := mocks.NewUserService(t)
	tokenServiceMock := mocks.NewTokenService(t)
	httpServer := NewHTTPServer(userServiceMock, tokenServiceMock, nil, nil, nil)

	req, err := http.NewRequestWithContext(context.Background(), "", "", nil)
	require.NoError(t, err)
	req.Header.Set(AuthorizationHeader, BearerPrefix+"invalid-token")

	tokenServiceMock.On("GetUser", "invalid-token").Return(domain.User{}, errors.New("invalid token"))

	rr := httptest.NewRecorder()
	handler := httpServer.CheckAuthorizedUser(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

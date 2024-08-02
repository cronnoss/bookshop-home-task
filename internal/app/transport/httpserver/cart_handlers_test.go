package httpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cronnoss/bookshop-home-task/internal/app/transport/httpserver/mocks"
	"github.com/stretchr/testify/assert"
)

func TestUpdateCart_InvalidUser(t *testing.T) {
	userServiceMock := mocks.NewUserService(t)
	cartServiceMock := mocks.NewCartService(t)

	httpServer := NewHTTPServer(userServiceMock, nil, nil, nil, cartServiceMock)

	reqBody := CartRequest{BookIDs: []int{1, 2}}
	reqBodyJSON, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/cart",
		bytes.NewBuffer(reqBodyJSON))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	httpServer.UpdateCart(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	userServiceMock.AssertNumberOfCalls(t, "GetUserByID", 0)
	cartServiceMock.AssertNumberOfCalls(t, "UpdateCartAndStocks", 0)
}

func TestUpdateCart_InvalidJSON(t *testing.T) {
	userServiceMock := mocks.NewUserService(t)
	cartServiceMock := mocks.NewCartService(t)

	httpServer := NewHTTPServer(userServiceMock, nil, nil, nil, cartServiceMock)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/cart",
		bytes.NewBuffer([]byte("invalid json")))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	httpServer.UpdateCart(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	userServiceMock.AssertNumberOfCalls(t, "GetUserByID", 0)
	cartServiceMock.AssertNumberOfCalls(t, "UpdateCartAndStocks", 0)
}

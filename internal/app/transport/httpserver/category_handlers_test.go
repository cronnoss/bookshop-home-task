package httpserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cronnoss/bookshop-home-task/internal/app/domain"
	"github.com/cronnoss/bookshop-home-task/internal/app/transport/httpserver/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetCategory_Success(t *testing.T) {
	categoryServiceMock := mocks.NewCategoryService(t)

	testCreatedCategory, err := domain.NewCategory(domain.NewCategoryData{
		Name: "Fiction",
	})
	require.NoError(t, err)

	categoryServiceMock.On("CreateCategory", mock.Anything, mock.Anything).Return(testCreatedCategory, nil).Once()

	httpServer := NewHTTPServer(nil, nil, nil, categoryServiceMock, nil)

	newCategoryRequest := []byte(`{
		  "name": "Fiction"
		}
	`)

	req := httptest.NewRequest(http.MethodPost, "/category", bytes.NewBuffer(newCategoryRequest))
	w := httptest.NewRecorder()

	httpServer.CreateCategory(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	// read response body
	var createCategoryResponse CategoryResponse
	err = json.NewDecoder(res.Body).Decode(&createCategoryResponse)
	assert.NoError(t, err)

	assert.Equal(t, createCategoryResponse.ID, testCreatedCategory.ID())
	assert.Equal(t, createCategoryResponse.Name, testCreatedCategory.Name())
}

func TestGetCategory_InvalidID(t *testing.T) {
	categoryServiceMock := mocks.NewCategoryService(t)
	httpServer := NewHTTPServer(nil, nil, nil, categoryServiceMock, nil)

	req := httptest.NewRequest(http.MethodGet, "/category/invalid", nil)
	w := httptest.NewRecorder()

	httpServer.GetCategory(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestCreateCategory_InvalidJSON(t *testing.T) {
	categoryServiceMock := mocks.NewCategoryService(t)
	httpServer := NewHTTPServer(nil, nil, nil, categoryServiceMock, nil)

	req := httptest.NewRequest(http.MethodPost, "/category", bytes.NewBuffer([]byte("invalid json")))
	w := httptest.NewRecorder()

	httpServer.CreateCategory(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestCreateCategory_InvalidRequest(t *testing.T) {
	categoryServiceMock := mocks.NewCategoryService(t)
	httpServer := NewHTTPServer(nil, nil, nil, categoryServiceMock, nil)

	invalidCategoryRequest := []byte(`{
		"name": ""
	}`)

	req := httptest.NewRequest(http.MethodPost, "/category", bytes.NewBuffer(invalidCategoryRequest))
	w := httptest.NewRecorder()

	httpServer.CreateCategory(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestUpdateCategory_ReturnsBadRequestForInvalidID(t *testing.T) {
	categoryServiceMock := mocks.NewCategoryService(t)
	httpServer := NewHTTPServer(nil, nil, nil, categoryServiceMock, nil)

	req := httptest.NewRequest(http.MethodPatch, "/category/invalid", nil)
	w := httptest.NewRecorder()

	httpServer.UpdateCategory(w, req)

	res := w.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestDeleteCategory_ReturnsBadRequestForInvalidID(t *testing.T) {
	categoryServiceMock := mocks.NewCategoryService(t)
	httpServer := NewHTTPServer(nil, nil, nil, categoryServiceMock, nil)

	req := httptest.NewRequest(http.MethodDelete, "/category/invalid", nil)
	w := httptest.NewRecorder()

	httpServer.DeleteCategory(w, req)

	res := w.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusBadRequest, res.StatusCode)
}

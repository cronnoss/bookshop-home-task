package httpserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cronnoss/bookshop-home-task/internal/app/domain"
	"github.com/cronnoss/bookshop-home-task/internal/app/transport/httpserver/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHttpServer_GetBook(t *testing.T) {
	bookServiceMock := mocks.NewBookService(t)

	testCreatedBook, err := domain.NewBook(domain.NewBookData{
		Title:      "The history of Golang",
		Year:       2024,
		Author:     "Rob Pike",
		Price:      1000,
		Stock:      100,
		CategoryID: 1,
	})
	require.NoError(t, err)

	bookServiceMock.On("CreateBook", mock.Anything, mock.Anything).Return(testCreatedBook, nil)

	httpServer := NewHTTPServer(nil, nil, bookServiceMock, nil, nil)

	newBookRequest := []byte(`{
		  "title": "The history of Golang",
		  "year": 2024,
		  "author": "Rob Pike",
		  "price": 1000,
		  "stock": 100,
		  "categoryId": 1
		}
	`)

	req := httptest.NewRequest(http.MethodPost, "/book", bytes.NewBuffer(newBookRequest))
	w := httptest.NewRecorder()

	httpServer.CreateBook(w, req)

	res := w.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)

	// read response body
	var createBookResponse BookResponse
	err = json.NewDecoder(res.Body).Decode(&createBookResponse)
	require.NoError(t, err)

	require.Equal(t, createBookResponse.ID, testCreatedBook.ID())
	require.Equal(t, createBookResponse.Title, testCreatedBook.Title())
	require.Equal(t, createBookResponse.Year, testCreatedBook.Year())
	require.Equal(t, createBookResponse.Author, testCreatedBook.Author())
	require.Equal(t, createBookResponse.Price, testCreatedBook.Price())
	require.Equal(t, createBookResponse.Stock, testCreatedBook.Stock())
	require.Equal(t, createBookResponse.CategoryID, testCreatedBook.CategoryID())
}

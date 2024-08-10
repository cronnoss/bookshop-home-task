package httpserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/cronnoss/bookshop-home-task/internal/app/common/server"
	"github.com/cronnoss/bookshop-home-task/internal/app/domain"
	"github.com/gorilla/mux"
)

// @Summary GetBook
// @Tags book
// @Description get book by ID
// @ID get-book
// @Accept  json
// @Produce  json
// @Param book_id path int true "book ID"
// @Success 200 {object} BookResponse
// @Failure 400,404 {object} server.ErrorResponse
// @Router /book/{book_id} [get]
// GetBook returns a book by ID
func (h HTTPServer) GetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID, err := strconv.Atoi(vars["book_id"])
	if err != nil {
		server.BadRequest("invalid-book-id", err, w, r)
		return
	}
	book, err := h.bookService.GetBook(r.Context(), bookID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.NotFound("book-not-found", err, w, r)
			return
		}
		server.RespondWithError(err, w, r)
		return
	}

	response := toResponseBook(book)

	server.RespondOK(response, w, r)
}

// @Summary CreateBook
// @Security ApiKeyAuth
// @Tags book
// @Description create book
// @ID create-book
// @Accept  json
// @Produce  json
// @Param input body BookRequest true "book info"
// @Success 200 {object} BookResponse
// @Failure 400,404 {object} server.ErrorResponse
// @Failure 401 {object} server.ErrorResponse
// @Failure 500 {object} server.ErrorResponse
// @Router /book [post]
// CreateBook creates a new book
func (h HTTPServer) CreateBook(w http.ResponseWriter, r *http.Request) {
	var bookRequest BookRequest
	if err := json.NewDecoder(r.Body).Decode(&bookRequest); err != nil {
		server.BadRequest("invalid-json", err, w, r)
		return
	}

	if err := bookRequest.Validate(); err != nil {
		server.BadRequest("invalid-request", err, w, r)
		return
	}

	book, err := toDomainBook(bookRequest)
	if err != nil {
		server.RespondWithError(err, w, r)
		return
	}

	insertedBook, err := h.bookService.CreateBook(r.Context(), book)
	if err != nil {
		server.RespondWithError(err, w, r)
		return
	}

	response := toResponseBook(insertedBook)

	server.RespondOK(response, w, r)
}

// @Summary UpdateBook
// @Security ApiKeyAuth
// @Tags book
// @Description update book by ID
// @ID update-book
// @Accept  json
// @Produce  json
// @Param book_id path int true "book ID"
// @Param input body BookRequest true "book info"
// @Success 200 {object} BookResponse
// @Failure 400,404 {object} server.ErrorResponse
// @Failure 401 {object} server.ErrorResponse
// @Failure 500 {object} server.ErrorResponse
// @Router /book/{book_id} [patch]
// UpdateBook updates a book by ID
func (h HTTPServer) UpdateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID, err := strconv.Atoi(vars["book_id"])
	if err != nil {
		server.BadRequest("invalid-book-id", err, w, r)
		return
	}

	var bookRequest BookRequest
	if err := json.NewDecoder(r.Body).Decode(&bookRequest); err != nil {
		server.BadRequest("invalid-json", err, w, r)
		return
	}

	if err := bookRequest.Validate(); err != nil {
		server.BadRequest("invalid-request", err, w, r)
		return
	}

	_, err = h.bookService.GetBook(r.Context(), bookID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.NotFound("book-not-found", err, w, r)
			return
		}
		server.RespondWithError(err, w, r)
		return
	}

	book, err := domain.NewBook(domain.NewBookData{
		ID:         bookID,
		Title:      bookRequest.Title,
		Year:       bookRequest.Year,
		Author:     bookRequest.Author,
		Price:      bookRequest.Price,
		CategoryID: bookRequest.CategoryID,
	})
	if err != nil {
		server.RespondWithError(err, w, r)
		return
	}

	updatedBook, err := h.bookService.UpdateBook(r.Context(), book)
	if err != nil {
		server.RespondWithError(err, w, r)
		return
	}

	response := toResponseBook(updatedBook)

	server.RespondOK(response, w, r)
}

// @Summary DeleteBook
// @Security ApiKeyAuth
// @Tags book
// @Description delete book by ID
// @ID delete-book
// @Accept  json
// @Produce  json
// @Param book_id path int true "book ID"
// @Success 200 {object} map[string]bool
// @Failure 400,404 {object} server.ErrorResponse
// @Failure 401 {object} server.ErrorResponse
// @Failure 500 {object} server.ErrorResponse
// @Router /book/{book_id} [delete]
// DeleteBook deletes a book by ID
func (h HTTPServer) DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID, err := strconv.Atoi(vars["book_id"])
	if err != nil {
		server.BadRequest("invalid-book-id", err, w, r)
		return
	}

	_, err = h.bookService.GetBook(r.Context(), bookID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.NotFound("book-not-found", err, w, r)
			return
		}
		server.RespondWithError(err, w, r)
		return
	}

	err = h.bookService.DeleteBook(r.Context(), bookID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			server.NotFound("book-not-found", err, w, r)
			return
		}
		server.RespondWithError(err, w, r)
		return
	}

	server.RespondOK(map[string]bool{"deleted": true}, w, r)
}

// @Summary GetBooks
// @Tags book
// @Description get books
// @ID get-books
// @Accept  json
// @Produce  json
// @Param category_id query []int false "category ID"
// @Param page query int false "page number"
// @Success 200 {array} BookResponse
// @Failure 400,404 {object} server.ErrorResponse
// @Router /books [get]
func (h HTTPServer) GetBooks(w http.ResponseWriter, r *http.Request) {
	// filter by category IDs
	queryCategoryIDs := r.URL.Query()["category_id"]
	categoryIDs := make([]int, 0, len(queryCategoryIDs))
	for _, id := range queryCategoryIDs {
		categoryID, err := strconv.Atoi(id)
		if err != nil {
			server.BadRequest("invalid-category-id", err, w, r)
			return
		}
		categoryIDs = append(categoryIDs, categoryID)
	}
	// page
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	var limit, offset int
	if page > 0 {
		limit = 10
		offset = (page - 1) * limit
	}

	books, err := h.bookService.GetBooks(r.Context(), categoryIDs, limit, offset)
	if err != nil {
		server.RespondWithError(err, w, r)
		return
	}

	response := make([]BookResponse, 0, len(books))
	for _, book := range books {
		response = append(response, toResponseBook(book))
	}

	server.RespondOK(response, w, r)
}

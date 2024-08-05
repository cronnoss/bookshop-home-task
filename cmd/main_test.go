package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cronnoss/bookshop-home-task/internal/app/config"
	"github.com/cronnoss/bookshop-home-task/internal/app/repository/pgrepo"
	"github.com/cronnoss/bookshop-home-task/internal/app/services"
	"github.com/cronnoss/bookshop-home-task/internal/app/transport/httpserver"
	"github.com/cronnoss/bookshop-home-task/internal/pkg/pg"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestRunPgMigrations(t *testing.T) {
	t.Run("no migrations path provided", func(t *testing.T) {
		err := runPgMigrations("postgres://user:password@localhost:5432/dbname?sslmode=disable", "")
		assert.Error(t, err)
		assert.Equal(t, "no migrations path provided", err.Error())
	})

	t.Run("no DSN provided", func(t *testing.T) {
		err := runPgMigrations("", "file://migrations")
		assert.Error(t, err)
		assert.Equal(t, "no DSN provided", err.Error())
	})
}

func TestHTTPServer(t *testing.T) {
	cfg := config.Config{
		HTTPAddr:       ":8080",
		DSN:            "postgres://user:password@localhost:5432/dbname?sslmode=disable",
		MigrationsPath: "file://migrations",
	}

	pgDB, _ := pg.Dial(cfg.DSN)
	userRepo := pgrepo.NewUserRepo(pgDB)
	bookRepo := pgrepo.NewBookRepo(pgDB)
	categoryRepo := pgrepo.NewCategoryRepo(pgDB)
	cartRepo := pgrepo.NewCartRepo(pgDB)

	userService := services.NewUserService(userRepo)
	bookService := services.NewBookService(bookRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	tokenService := services.NewTokenService(tokenTTL)
	cartService := services.NewCartService(cartRepo)

	httpServer := httpserver.NewHTTPServer(userService, tokenService, bookService, categoryService, cartService)

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Book Shop API v0.1"))
	}).Methods("GET")

	router.HandleFunc("/signup", httpServer.SignUp).Methods(http.MethodPost)
	router.HandleFunc("/signin", httpServer.SignIn).Methods(http.MethodPost)

	router.HandleFunc("/books", httpServer.GetBooks).Methods(http.MethodGet)
	router.HandleFunc("/book/{book_id}", httpServer.GetBook).Methods(http.MethodGet)
	router.HandleFunc("/book", httpServer.CheckAdmin(httpServer.CreateBook)).Methods(http.MethodPost)
	router.HandleFunc("/book/{book_id}", httpServer.CheckAdmin(httpServer.UpdateBook)).Methods(http.MethodPatch)
	router.HandleFunc("/book/{book_id}", httpServer.CheckAdmin(httpServer.DeleteBook)).Methods(http.MethodDelete)

	router.HandleFunc("/categories", httpServer.GetCategories).Methods(http.MethodGet)
	router.HandleFunc("/category/{category_id}", httpServer.GetCategory).Methods(http.MethodGet)
	router.HandleFunc("/category", httpServer.CheckAdmin(httpServer.CreateCategory)).Methods(http.MethodPost)
	router.HandleFunc("/category/{category_id}", httpServer.CheckAdmin(httpServer.UpdateCategory)).
		Methods(http.MethodPatch)
	router.HandleFunc("/category/{category_id}", httpServer.CheckAdmin(httpServer.DeleteCategory)).
		Methods(http.MethodDelete)

	router.HandleFunc("/cart", httpServer.CheckAuthorizedUser(httpServer.UpdateCart)).Methods(http.MethodPost)
	router.HandleFunc("/checkout", httpServer.CheckAuthorizedUser(httpServer.Checkout)).Methods(http.MethodPost)

	t.Run("root endpoint", func(t *testing.T) {
		req, _ := http.NewRequestWithContext(context.Background(), "GET", "/", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "Book Shop API v0.1", rr.Body.String())
	})
}

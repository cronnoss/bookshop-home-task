// //go:build integration

package integration

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cronnoss/bookshop-home-task/internal/app/repository/models"
	"github.com/cronnoss/bookshop-home-task/internal/app/repository/pgrepo"
	servise "github.com/cronnoss/bookshop-home-task/internal/app/services"
	"github.com/cronnoss/bookshop-home-task/internal/app/transport/httpserver"
	"github.com/cronnoss/bookshop-home-task/internal/pkg/pg"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

const (
	dbUser1     = "user"
	dbPassword1 = "password"
)

type HTTPTestSuite struct {
	t               *testing.T
	userService     servise.UserService
	tokenService    servise.TokenService
	bookService     servise.BookService
	categoryService servise.CategoryService
	cartService     servise.CartService
	httpServer      httpserver.HTTPServer
	dbContainer     tc.Container
	db              *bun.DB
}

func (s *HTTPTestSuite) SetupSuite1() {
	ctx := context.Background()
	s.dbContainer = startPostgresContainer1(ctx, dbUser1, dbPassword1)
}

func (s *HTTPTestSuite) TearDownSuite1() {
	ctx := context.Background()
	ensureContainerTermination1(ctx, s.dbContainer)
}

func (s *HTTPTestSuite) TearDownTest1() {
	_, err := s.db.NewTruncateTable().Model((*models.User)(nil)).Exec(context.Background())
	require.NoError(s.t, err)
	_, err = s.db.NewTruncateTable().Model((*models.Book)(nil)).Exec(context.Background())
	require.NoError(s.t, err)
	_, err = s.db.NewTruncateTable().Model((*models.Category)(nil)).Exec(context.Background())
	require.NoError(s.t, err)
	_, err = s.db.NewTruncateTable().Model((*models.Cart)(nil)).Exec(context.Background())
	require.NoError(s.t, err)
}

func (s *HTTPTestSuite) prepareTestPostgresDatabase1(dbname string) *bun.DB {
	ctx := context.Background()

	// Initialize database
	initializeDatabase1(ctx, s.t, s.dbContainer, dbname, dbUser1, dbPassword1)

	mappedPort, err := s.dbContainer.MappedPort(ctx, "5432/tcp")
	require.NoError(s.t, err)

	// Establish database connection
	pg := connectToPostgresForTest1(s.t, "localhost", dbUser1, dbPassword1, dbname, mappedPort.Port())
	require.NoError(s.t, err)

	return pg
}

func TestHttpSuite(t *testing.T) {
	suite := &HTTPTestSuite{t: t}
	suite.SetupSuite1()
	defer suite.TearDownSuite1()

	t.Run("HttpSuite", func(t *testing.T) {
		t.Run("TestCreateBookHttp_Success", suite.TestCreateBookHttp_Success)
	})
}

func (s *HTTPTestSuite) TestCreateBookHttp_Success(t *testing.T) {
	s.db = s.prepareTestPostgresDatabase1(uuid.NewString())

	s.userService = servise.NewUserService(pgrepo.NewUserRepo(&pg.DB{DB: s.db}))
	s.tokenService = servise.NewTokenService(15)
	s.bookService = servise.NewBookService(pgrepo.NewBookRepo(&pg.DB{DB: s.db}))
	s.categoryService = servise.NewCategoryService(pgrepo.NewCategoryRepo(&pg.DB{DB: s.db}))
	s.cartService = servise.NewCartService(pgrepo.NewCartRepo(&pg.DB{DB: s.db}))

	// create http server with application injected
	s.httpServer = httpserver.NewHTTPServer(
		s.userService,
		s.tokenService,
		s.bookService,
		s.categoryService,
		s.cartService,
	)

	// 1. create POST /signup request
	newUserRequest := []byte(`{
		"username": "testuser",
		"password": "password123"
	}`)

	// create http request
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(newUserRequest))

	// create http recorder
	w := httptest.NewRecorder()

	// run request
	s.httpServer.SignUp(w, req)

	// get response
	res := w.Result()
	defer res.Body.Close()

	// read response body
	data, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	// implement response parsing.
	var response map[string]bool
	require.NoError(t, err)

	err = json.Unmarshal(data, &response) // { "ok": true }
	require.NoError(t, err)

	_, ok := response["ok"]
	require.True(t, ok)
	assert.Equal(t, true, response["ok"])

	require.Equal(t, http.StatusOK, res.StatusCode)

	// 2. create POST /signin request
	signInRequest := []byte(`{
		"username": "testuser",
		"password": "password123"
	}`)

	// create http request
	req = httptest.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(signInRequest))

	// create http recorder
	w = httptest.NewRecorder()

	// run request
	s.httpServer.SignIn(w, req)

	// get response
	res = w.Result()
	defer res.Body.Close()

	// read response body
	data, err = io.ReadAll(res.Body)
	require.NoError(t, err)

	// implement token parsing
	var response1 map[string]string

	err = json.Unmarshal(data, &response1)
	require.NoError(t, err)

	token, ok := response1["token"]
	require.True(t, ok)
	assert.NotEmpty(t, token)

	require.Equal(t, http.StatusOK, res.StatusCode)

	// 3. create POST /book request
	newBookRequest := []byte(`{
		"title": "1984",
		"year": 1949,
		"author": "George Orwell",
		"price": 1500,
		"stock": 200,
		"categoryId": 1
	}`)

	// create http request with token for authorization
	req = httptest.NewRequest(http.MethodPost, "/book", bytes.NewBuffer(newBookRequest))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// create http recorder
	w = httptest.NewRecorder()

	// run request
	s.httpServer.CreateBook(w, req)

	// get response
	res = w.Result()
	defer res.Body.Close()

	// read response body
	data, err = io.ReadAll(res.Body)
	require.NoError(t, err)

	// unmarshal response
	var responseBook models.Book
	err = json.Unmarshal(data, &responseBook)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "1984", responseBook.Title)
	require.Equal(t, 1949, responseBook.Year)
	require.Equal(t, "George Orwell", responseBook.Author)
	require.Equal(t, 1500, responseBook.Price)
	require.Equal(t, 200, responseBook.Stock)
	require.Equal(t, 1, responseBook.CategoryID)
}

func connectToPostgresForTest1(_ *testing.T, host, user, password, dbname, port string) *bun.DB {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	if err := createSchema1(context.Background(), db); err != nil {
		log.Fatalf("Schema creation failed: %v", err)
	}
	log.Print("Successfully connected to database")

	return db
}

func createSchema1(ctx context.Context, db *bun.DB) error {
	_, err := db.NewCreateTable().Model((*models.User)(nil)).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create user table: %w", err)
	}
	_, err = db.NewCreateTable().Model((*models.Book)(nil)).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create books table: %w", err)
	}
	_, err = db.NewCreateTable().Model((*models.Category)(nil)).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create categories table: %w", err)
	}
	_, err = db.NewCreateTable().Model((*models.Cart)(nil)).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create carts table: %w", err)
	}
	_, err = db.ExecContext(ctx, `ALTER TABLE carts ADD CONSTRAINT unique_user_id UNIQUE (user_id)`)
	if err != nil {
		return fmt.Errorf("failed to add unique constraint to carts table: %w", err)
	}
	return nil
}

func initializeDatabase1(
	ctx context.Context,
	t *testing.T,
	container tc.Container,
	dbname, user, _ string,
) {
	t.Helper()

	exitCode, _, err := container.Exec(ctx, []string{
		"createdb",
		"-p", "5432",
		"-h", "localhost",
		"-U", user,
		dbname,
	})

	require.NoError(t, err)
	require.Equal(t, 0, exitCode, "Non-zero exit code from 'createdb'")
}

func startPostgresContainer1(ctx context.Context, user, password string) tc.Container {
	req := tc.ContainerRequest{
		Image: "postgres:latest",
		Env: map[string]string{
			"POSTGRES_USER":     user,
			"POSTGRES_PASSWORD": password,
		},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
	}
	container, err := tc.GenericContainer(
		ctx,
		tc.GenericContainerRequest{ContainerRequest: req, Started: true},
	)
	if err != nil {
		log.Fatal("Failed to start test container")
	}

	return container
}

func ensureContainerTermination1(ctx context.Context, container tc.Container) {
	if err := container.Terminate(ctx); err != nil {
		log.Fatalf("Failed to terminate container: %v", err)
	}
}

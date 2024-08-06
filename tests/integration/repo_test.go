// //go:build integration

package integration

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/cronnoss/bookshop-home-task/internal/app/domain"
	"github.com/cronnoss/bookshop-home-task/internal/app/repository/models"
	"github.com/cronnoss/bookshop-home-task/internal/app/repository/pgrepo"
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
	dbUser     = "user"
	dbPassword = "password"
)

type Book struct {
	Title      string
	Year       int
	Author     string
	Price      int
	Stock      int
	CategoryID int
}

type IntegrationSuite struct {
	t           *testing.T
	dbContainer tc.Container
	db          *bun.DB
}

type PGDBAdapter struct {
	*bun.DB
}

func NewPGDBAdapter(db *bun.DB) *PGDBAdapter {
	return &PGDBAdapter{db}
}

func (s *IntegrationSuite) SetupSuite() {
	ctx := context.Background()
	s.dbContainer = startPostgresContainer(ctx, dbUser, dbPassword)
}

func (s *IntegrationSuite) TearDownSuite() {
	ctx := context.Background()
	ensureContainerTermination(ctx, s.dbContainer)
}

func (s *IntegrationSuite) TearDownTest() {
	_, err := s.db.NewTruncateTable().Model((*models.User)(nil)).Exec(context.Background())
	require.NoError(s.t, err)
}

func (s *IntegrationSuite) prepareTestPostgresDatabase(dbname string) *bun.DB {
	ctx := context.Background()

	// Initialize database
	initializeDatabase(ctx, s.t, s.dbContainer, dbname, dbUser, dbPassword)

	mappedPort, err := s.dbContainer.MappedPort(ctx, "5432/tcp")
	require.NoError(s.t, err)

	// Establish database connection
	pg := connectToPostgresForTest(
		s.t,
		"localhost", dbUser, dbPassword, dbname, mappedPort.Port(),
	)
	require.NoError(s.t, err)

	return pg
}

func TestIntegrationSuite(t *testing.T) {
	suite := &IntegrationSuite{t: t}
	suite.SetupSuite()
	defer suite.TearDownSuite()

	t.Run("IntegrationSuite", func(t *testing.T) {
		// UserRepo tests
		t.Run("TestCreateUser_Success", suite.TestCreateUser_Success)
		t.Run("TestCreateUser_FailsOnInsert", suite.TestCreateUser_FailsOnInsert)
		t.Run("TestGetUser_NotFound", suite.TestGetUser_NotFound)
		// BookRepo tests
		t.Run("TestCreateBook_Success", suite.TestCreateBook_Success)
		t.Run("TestGetBook_Success", suite.TestGetBook_Success)
		t.Run("TestGetBook_NotFound", suite.TestGetBook_NotFound)
		t.Run("TestUpdateBook_Success", suite.TestUpdateBook_Success)
		t.Run("TestDeleteBook_Success", suite.TestDeleteBook_Success)
		t.Run("TestGetBooks_Success", suite.TestGetBooks_Success)
		// CartRepo tests
		t.Run("TestGetCart_Success", suite.TestGetCart_Success)
		t.Run("TestGetCart_NotFound", suite.TestGetCart_NotFound)
		t.Run("TestUpdateCartAndStocks_Success", suite.TestUpdateCartAndStocks_Success)
		t.Run("TestCheckStocks_Success", suite.TestCheckStocks_Success)
		t.Run("TestDeleteCart_Success", suite.TestDeleteCart_Success)
		// HandleBunTransaction tests
		t.Run("TestHandleBunTransaction_Success", suite.TestHandleBunTransaction_Success)
		t.Run("TestHandleBunTransaction_FailBegin", suite.TestHandleBunTransaction_FailBegin)
		t.Run("TestHandleBunTransaction_FailCommit", suite.TestHandleBunTransaction_FailCommit)
		t.Run("TestHandleBunTransaction_FailRollback", suite.TestHandleBunTransaction_FailRollback)
	})
}

// UserRepo tests.
func (s *IntegrationSuite) TestCreateUser_Success(t *testing.T) {
	ctx := context.Background()

	s.db = s.prepareTestPostgresDatabase(uuid.NewString())

	userRepo := pgrepo.UserRepo{DB: (*pg.DB)(NewPGDBAdapter(s.db))}
	user := domain.User{ID: 1, Username: "testuser", Password: "password123"}

	createdUser, err := userRepo.CreateUser(ctx, user)
	require.NoError(t, err)

	assert.Equal(t, "testuser", createdUser.Username)
	assert.Equal(t, "password123", createdUser.Password)

	// Retrieve the user from the database
	retrievedUser, err := userRepo.GetUser(ctx, "testuser")
	require.NoError(t, err)
	assert.Equal(t, "testuser", retrievedUser.Username)
	assert.Equal(t, "password123", retrievedUser.Password)
}

func (s *IntegrationSuite) TestCreateUser_FailsOnInsert(t *testing.T) {
	ctx := context.Background()

	s.db = s.prepareTestPostgresDatabase(uuid.NewString())

	userRepo := pgrepo.UserRepo{DB: (*pg.DB)(NewPGDBAdapter(s.db))}
	user := domain.User{ID: 1, Username: "testuser", Password: "password123"}

	_, err := userRepo.CreateUser(ctx, user)
	require.NoError(t, err)

	_, err = userRepo.CreateUser(ctx, user)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate key value violates unique constraint \"users_pkey\"")
}

func (s *IntegrationSuite) TestGetUser_NotFound(t *testing.T) {
	ctx := context.Background()

	s.db = s.prepareTestPostgresDatabase(uuid.NewString())

	userRepo := pgrepo.UserRepo{DB: (*pg.DB)(NewPGDBAdapter(s.db))}

	// Attempt to retrieve a user that does not exist
	_, err := userRepo.GetUser(ctx, "testuser")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// BookRepo tests.
func (s *IntegrationSuite) TestCreateBook_Success(t *testing.T) {
	ctx := context.Background()

	s.db = s.prepareTestPostgresDatabase(uuid.NewString())

	bookRepo := pgrepo.NewBookRepo(&pg.DB{DB: s.db})

	bookData := domain.NewBookData{
		Title:      "1984",
		Year:       1949,
		Author:     "George Orwell",
		Price:      1500,
		Stock:      200,
		CategoryID: 1,
	}

	book, err := domain.NewBook(bookData)
	require.NoError(t, err)

	createdBook, err := bookRepo.CreateBook(ctx, book)
	require.NoError(t, err)

	assert.Equal(t, "1984", createdBook.Title())
	assert.Equal(t, "George Orwell", createdBook.Author())
	assert.Equal(t, 1949, createdBook.Year())
	assert.Equal(t, 1500, createdBook.Price())
	assert.Equal(t, 200, createdBook.Stock())
	assert.Equal(t, 1, createdBook.CategoryID())
}

func (s *IntegrationSuite) TestGetBook_Success(t *testing.T) {
	ctx := context.Background()

	s.db = s.prepareTestPostgresDatabase(uuid.NewString())

	bookRepo := pgrepo.NewBookRepo(&pg.DB{DB: s.db})

	bookData := domain.NewBookData{
		Title:      "1984",
		Year:       1949,
		Author:     "George Orwell",
		Price:      1500,
		Stock:      200,
		CategoryID: 1,
	}

	book, err := domain.NewBook(bookData)
	require.NoError(t, err)

	createdBook, err := bookRepo.CreateBook(ctx, book)
	require.NoError(t, err)

	retrievedBook, err := bookRepo.GetBook(ctx, createdBook.ID())
	require.NoError(t, err)

	assert.Equal(t, "1984", retrievedBook.Title())
	assert.Equal(t, "George Orwell", retrievedBook.Author())
	assert.Equal(t, 1949, retrievedBook.Year())
	assert.Equal(t, 1500, retrievedBook.Price())
	assert.Equal(t, 200, retrievedBook.Stock())
	assert.Equal(t, 1, retrievedBook.CategoryID())
}

func (s *IntegrationSuite) TestGetBook_NotFound(t *testing.T) {
	ctx := context.Background()

	s.db = s.prepareTestPostgresDatabase(uuid.NewString())

	bookRepo := pgrepo.NewBookRepo(&pg.DB{DB: s.db})

	_, err := bookRepo.GetBook(ctx, 1)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func (s *IntegrationSuite) TestUpdateBook_Success(t *testing.T) {
	ctx := context.Background()

	s.db = s.prepareTestPostgresDatabase(uuid.NewString())

	bookRepo := pgrepo.NewBookRepo(&pg.DB{DB: s.db})

	bookData := domain.NewBookData{
		Title:      "1984",
		Year:       1949,
		Author:     "George Orwell",
		Price:      1500,
		Stock:      200,
		CategoryID: 1,
	}

	book, err := domain.NewBook(bookData)
	require.NoError(t, err)

	createdBook, err := bookRepo.CreateBook(ctx, book)
	require.NoError(t, err)

	createdBookData := domain.NewBookData{
		ID:         createdBook.ID(),
		Title:      "Animal Farm",
		Year:       1945,
		Author:     "George Orwell",
		Price:      1000,
		Stock:      100,
		CategoryID: 1,
	}

	updatedBook, err := domain.NewBook(createdBookData)
	require.NoError(t, err)

	updatedBook, err = bookRepo.UpdateBook(ctx, updatedBook)
	require.NoError(t, err)

	assert.Equal(t, "Animal Farm", updatedBook.Title())
	assert.Equal(t, "George Orwell", updatedBook.Author())
	assert.Equal(t, 1945, updatedBook.Year())
	assert.Equal(t, 1000, updatedBook.Price())
	assert.Equal(t, 200, updatedBook.Stock()) // Stock should not be updated
	assert.Equal(t, 1, updatedBook.CategoryID())
}

func (s *IntegrationSuite) TestDeleteBook_Success(t *testing.T) {
	ctx := context.Background()

	s.db = s.prepareTestPostgresDatabase(uuid.NewString())

	bookRepo := pgrepo.NewBookRepo(&pg.DB{DB: s.db})

	bookData := domain.NewBookData{
		Title:      "1984",
		Year:       1949,
		Author:     "George Orwell",
		Price:      1500,
		Stock:      200,
		CategoryID: 1,
	}

	book, err := domain.NewBook(bookData)
	require.NoError(t, err)

	createdBook, err := bookRepo.CreateBook(ctx, book)
	require.NoError(t, err)

	err = bookRepo.DeleteBook(ctx, createdBook.ID())
	require.NoError(t, err)

	_, err = bookRepo.GetBook(ctx, createdBook.ID())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func (s *IntegrationSuite) TestGetBooks_Success(t *testing.T) {
	ctx := context.Background()

	s.db = s.prepareTestPostgresDatabase(uuid.NewString())

	bookRepo := pgrepo.NewBookRepo(&pg.DB{DB: s.db})

	bookData1 := domain.NewBookData{
		Title:      "1984",
		Year:       1949,
		Author:     "George Orwell",
		Price:      1500,
		Stock:      200,
		CategoryID: 1,
	}

	bookData2 := domain.NewBookData{
		Title:      "Animal Farm",
		Year:       1945,
		Author:     "George Orwell",
		Price:      1000,
		Stock:      100,
		CategoryID: 1,
	}

	book1, err := domain.NewBook(bookData1)
	require.NoError(t, err)

	book2, err := domain.NewBook(bookData2)
	require.NoError(t, err)

	_, err = bookRepo.CreateBook(ctx, book1)
	require.NoError(t, err)

	_, err = bookRepo.CreateBook(ctx, book2)
	require.NoError(t, err)

	books, err := bookRepo.GetBooks(ctx, []int{1}, 10, 0)
	require.NoError(t, err)

	assert.Len(t, books, 2)
}

// CartRepo tests.
func (s *IntegrationSuite) TestGetCart_Success(t *testing.T) {
	ctx := context.Background()

	s.db = s.prepareTestPostgresDatabase(uuid.NewString())

	cartRepo := pgrepo.NewCartRepo(&pg.DB{DB: s.db})

	cartData := domain.NewCartData{
		UserID:  1,
		BookIDs: []int{1, 2},
	}

	cart, err := domain.NewCart(cartData)
	require.NoError(t, err)
	cartModel := &models.Cart{
		UserID:  cart.UserID(),
		BookIDs: cart.BookIDs(),
	}
	_, err = s.db.NewInsert().Model(cartModel).Column("user_id", "book_ids").Exec(ctx)
	require.NoError(t, err)

	retrievedCart, err := cartRepo.GetCart(ctx, cartData.UserID)
	require.NoError(t, err)

	assert.Equal(t, cartData.UserID, retrievedCart.UserID())
	assert.ElementsMatch(t, cartData.BookIDs, retrievedCart.BookIDs())
}

func (s *IntegrationSuite) TestGetCart_NotFound(t *testing.T) {
	ctx := context.Background()

	s.db = s.prepareTestPostgresDatabase(uuid.NewString())

	cartRepo := pgrepo.NewCartRepo(&pg.DB{DB: s.db})

	_, err := cartRepo.GetCart(ctx, 1)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func (s *IntegrationSuite) TestUpdateCartAndStocks_Success(t *testing.T) {
	ctx := context.Background()

	s.db = s.prepareTestPostgresDatabase(uuid.NewString())

	cartRepo := pgrepo.NewCartRepo(&pg.DB{DB: s.db})

	cartData := domain.NewCartData{
		UserID:  1,
		BookIDs: []int{1, 2},
	}

	cart, err := domain.NewCart(cartData)
	require.NoError(t, err)
	cartModel := &models.Cart{
		UserID:  cart.UserID(),
		BookIDs: cart.BookIDs(),
	}
	_, err = s.db.NewInsert().Model(cartModel).Column("user_id", "book_ids").
		On("CONFLICT (user_id) DO UPDATE").Set("book_ids = EXCLUDED.book_ids").Exec(ctx)
	require.NoError(t, err)

	_, err = cartRepo.GetCart(ctx, cartData.UserID)
	require.NoError(t, err)

	updatedCartData := domain.NewCartData{
		UserID:  1,
		BookIDs: []int{1},
	}

	updatedCart, err := domain.NewCart(updatedCartData)
	require.NoError(t, err)
	_ = &models.Cart{
		UserID:  updatedCart.UserID(),
		BookIDs: updatedCart.BookIDs(),
	}

	err = cartRepo.UpdateCartAndStocks(ctx, updatedCart)
	require.NoError(t, err)

	retrievedCart, err := cartRepo.GetCart(ctx, updatedCartData.UserID)
	require.NoError(t, err)

	assert.Equal(t, updatedCartData.UserID, retrievedCart.UserID())
	assert.ElementsMatch(t, updatedCartData.BookIDs, retrievedCart.BookIDs())
}

func (s *IntegrationSuite) TestCheckStocks_Success(t *testing.T) {
	ctx := context.Background()

	s.db = s.prepareTestPostgresDatabase(uuid.NewString())

	cartRepo := pgrepo.NewCartRepo(&pg.DB{DB: s.db})

	bookData1 := domain.NewBookData{
		Title:      "1984",
		Year:       1949,
		Author:     "George Orwell",
		Price:      1500,
		Stock:      200,
		CategoryID: 1,
	}

	bookData2 := domain.NewBookData{
		Title:      "Animal Farm",
		Year:       1945,
		Author:     "George Orwell",
		Price:      1000,
		Stock:      100,
		CategoryID: 1,
	}

	book1, err := domain.NewBook(bookData1)
	require.NoError(t, err)

	book2, err := domain.NewBook(bookData2)
	require.NoError(t, err)

	_, err = pgrepo.NewBookRepo(&pg.DB{DB: s.db}).CreateBook(ctx, book1)
	require.NoError(t, err)

	_, err = pgrepo.NewBookRepo(&pg.DB{DB: s.db}).CreateBook(ctx, book2)
	require.NoError(t, err)

	cartData := domain.NewCartData{
		UserID:  1,
		BookIDs: []int{1, 2},
	}

	cart, err := domain.NewCart(cartData)
	require.NoError(t, err)

	err = cartRepo.UpdateCartAndStocks(ctx, cart)
	require.NoError(t, err)

	ok, err := cartRepo.CheckStocks(ctx, cart)
	require.NoError(t, err)
	assert.True(t, ok)

	updatedCart, err := cartRepo.GetCart(ctx, cartData.UserID)
	require.NoError(t, err)

	assert.Equal(t, cartData.UserID, updatedCart.UserID())

	// Check that the stock of the books has been updated
	bookRepo := pgrepo.NewBookRepo(&pg.DB{DB: s.db})

	book1, err = bookRepo.GetBook(ctx, 1)
	require.NoError(t, err)

	book2, err = bookRepo.GetBook(ctx, 2)
	require.NoError(t, err)

	assert.Equal(t, 199, book1.Stock())
	assert.Equal(t, 99, book2.Stock())
}

func (s *IntegrationSuite) TestDeleteCart_Success(t *testing.T) {
	ctx := context.Background()

	s.db = s.prepareTestPostgresDatabase(uuid.NewString())

	cartRepo := pgrepo.NewCartRepo(&pg.DB{DB: s.db})

	cartData := domain.NewCartData{
		UserID:  1,
		BookIDs: []int{1, 2},
	}

	cart, err := domain.NewCart(cartData)
	require.NoError(t, err)
	cartModel := &models.Cart{
		UserID:  cart.UserID(),
		BookIDs: cart.BookIDs(),
	}
	_, err = s.db.NewInsert().Model(cartModel).Column("user_id", "book_ids").Exec(ctx)
	require.NoError(t, err)

	err = cartRepo.DeleteCart(ctx, cartData.UserID)
	require.NoError(t, err)

	_, err = cartRepo.GetCart(ctx, cartData.UserID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// HandleBunTransaction tests.
func (s *IntegrationSuite) TestHandleBunTransaction_Success(t *testing.T) {
	ctx := context.Background()

	s.db = s.prepareTestPostgresDatabase(uuid.NewString())

	pgDB := NewPGDBAdapter(s.db)

	db1 := &pg.DB{DB: pgDB.DB}
	bunTx := func(tx bun.Tx) error {
		book := &Book{
			Title:      "The Great Gatsby",
			Year:       1925,
			Author:     "F. Scott Fitzgerald",
			Price:      1000,
			Stock:      100,
			CategoryID: 1,
		}

		_, err := tx.NewInsert().Model(book).Exec(ctx)
		if err != nil {
			return err
		}

		return nil
	}

	err := pg.HandleBunTransaction(ctx, bunTx, db1)
	assert.NoError(t, err)
}

func (s *IntegrationSuite) TestHandleBunTransaction_FailBegin(t *testing.T) {
	ctx := context.Background()

	s.db = s.prepareTestPostgresDatabase(uuid.NewString())

	pgDB := NewPGDBAdapter(s.db)

	db1 := &pg.DB{DB: pgDB.DB}
	bunTx := func(tx bun.Tx) error {
		return errors.New("begin transaction failed")
	}

	err := pg.HandleBunTransaction(ctx, bunTx, db1)
	assert.Error(t, err)
	assert.Equal(t, "failed executing transaction: begin transaction failed", err.Error())
}

func (s *IntegrationSuite) TestHandleBunTransaction_FailCommit(t *testing.T) {
	ctx := context.Background()

	s.db = s.prepareTestPostgresDatabase(uuid.NewString())
	pgDB := NewPGDBAdapter(s.db)

	db1 := &pg.DB{DB: pgDB.DB}
	bunTx := func(tx bun.Tx) error {
		return errors.New("commit transaction failed")
	}

	err := pg.HandleBunTransaction(ctx, bunTx, db1)
	assert.Error(t, err)
	assert.Equal(t, "failed executing transaction: commit transaction failed", err.Error())
}

func (s *IntegrationSuite) TestHandleBunTransaction_FailRollback(t *testing.T) {
	ctx := context.Background()

	s.db = s.prepareTestPostgresDatabase(uuid.NewString())
	pgDB := NewPGDBAdapter(s.db)

	db1 := &pg.DB{DB: pgDB.DB}
	bunTx := func(tx bun.Tx) error {
		return errors.New("rollback transaction failed")
	}

	err := pg.HandleBunTransaction(ctx, bunTx, db1)
	assert.Error(t, err)
	assert.Equal(t, "failed executing transaction: rollback transaction failed", err.Error())
}

func connectToPostgresForTest(_ *testing.T, host, user, password, dbname, port string) *bun.DB {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	if err := createSchema(context.Background(), db); err != nil {
		log.Fatalf("Schema creation failed: %v", err)
	}
	log.Print("Successfully connected to database")

	return db
}

func createSchema(ctx context.Context, db *bun.DB) error {
	_, err := db.NewCreateTable().Model((*models.User)(nil)).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create user table: %w", err)
	}
	_, err = db.NewCreateTable().Model((*models.Book)(nil)).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create books table: %w", err)
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

func initializeDatabase(
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

func startPostgresContainer(ctx context.Context, user, password string) tc.Container {
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

func ensureContainerTermination(ctx context.Context, container tc.Container) {
	if err := container.Terminate(ctx); err != nil {
		log.Fatalf("Failed to terminate container: %v", err)
	}
}

// //go:build integration

package pgrepo

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"

	"github.com/cronnoss/bookshop-home-task/internal/app/domain"
	"github.com/cronnoss/bookshop-home-task/internal/app/repository/models"
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

type IntegrationSuite struct {
	t           *testing.T
	dbContainer tc.Container
	db          *bun.DB
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
		t.Run("TestCreateUser_Success", suite.TestCreateUser_Success)
		t.Run("TestCreateUser_FailsOnInsert", suite.TestCreateUser_FailsOnInsert)
		t.Run("TestGetUser_NotFound", suite.TestGetUser_NotFound)
	})
}

func (s *IntegrationSuite) TestCreateUser_Success(t *testing.T) {
	ctx := context.Background()

	s.db = s.prepareTestPostgresDatabase(uuid.NewString())

	user := domain.User{ID: 1, Username: "testuser", Password: "password123"}

	_, err := s.db.NewInsert().Model(&user).Exec(ctx)
	require.NoError(t, err)

	err = s.db.NewSelect().
		Model(&user).
		Where("id = 1").
		Scan(ctx)

	require.NoError(t, err)
	assert.NotNil(t, user.Admin)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "password123", user.Password)
}

func (s *IntegrationSuite) TestCreateUser_FailsOnInsert(t *testing.T) {
	ctx := context.Background()

	s.db = s.prepareTestPostgresDatabase(uuid.NewString())

	user := domain.User{ID: 1, Username: "testuser", Password: "password123"}

	_, err := s.db.NewInsert().Model(&user).Exec(ctx)
	require.NoError(t, err)

	_, err = s.db.NewInsert().Model(&user).Exec(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate key value violates unique constraint \"users_pkey\"")
}

func (s *IntegrationSuite) TestGetUser_NotFound(t *testing.T) {
	ctx := context.Background()

	s.db = s.prepareTestPostgresDatabase(uuid.NewString())

	user := domain.User{ID: 1, Username: "testuser", Password: "password123"}

	err := s.db.NewSelect().
		Model(&user).
		Where("id = 1").
		Scan(ctx)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no rows in result set")
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

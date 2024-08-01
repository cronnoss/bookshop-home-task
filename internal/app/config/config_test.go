package config

import (
	"os"
	"testing"
)

func TestReadWithAllEnvVarsSet(t *testing.T) {
	os.Setenv("HTTP_ADDR", ":8080")
	os.Setenv("DSN", "postgres://postgres:password@127.0.0.1:5433/bookshop?sslmode=disable")
	os.Setenv("MIGRATIONS_PATH", "file://./internal/app/migrations")
	defer os.Clearenv()

	config := Read()

	if config.HTTPAddr != ":8080" {
		t.Errorf("expected HTTPAddr to be ':8080', got '%s'", config.HTTPAddr)
	}
	if config.DSN != "postgres://postgres:password@127.0.0.1:5433/bookshop?sslmode=disable" {
		t.Errorf("expected DSN to be 'postgres://postgres:password@127.0.0.1:5433/bookshop?sslmode=disable', "+
			"got '%s'", config.DSN)
	}
	if config.MigrationsPath != "file://./internal/app/migrations" {
		t.Errorf("expected MigrationsPath to be 'file://./internal/app/migrations', got '%s'", config.MigrationsPath)
	}
}

func TestReadWithNoEnvVarsSet(t *testing.T) {
	os.Clearenv()

	config := Read()

	if config.HTTPAddr != "" {
		t.Errorf("expected HTTPAddr to be empty, got '%s'", config.HTTPAddr)
	}
	if config.DSN != "" {
		t.Errorf("expected DSN to be empty, got '%s'", config.DSN)
	}
	if config.MigrationsPath != "" {
		t.Errorf("expected MigrationsPath to be empty, got '%s'", config.MigrationsPath)
	}
}

func TestReadWithPartialEnvVarsSet(t *testing.T) {
	os.Setenv("HTTP_ADDR", ":8080")
	defer os.Clearenv()

	config := Read()

	if config.HTTPAddr != ":8080" {
		t.Errorf("expected HTTPAddr to be ':8080', got '%s'", config.HTTPAddr)
	}
	if config.DSN != "" {
		t.Errorf("expected DSN to be empty, got '%s'", config.DSN)
	}
	if config.MigrationsPath != "" {
		t.Errorf("expected MigrationsPath to be empty, got '%s'", config.MigrationsPath)
	}
}

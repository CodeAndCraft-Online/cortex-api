package testutils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/CodeAndCraft-Online/cortex-api/internal/database"
	"github.com/CodeAndCraft-Online/cortex-api/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupMockDB creates a sqlmock database for unit testing
func SetupMockDB() (*sql.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	return db, mock, nil
}

// SetupTestDB creates a real test database using dockertest (local) or CI service
func SetupTestDB() (*gorm.DB, func(), error) {
	if os.Getenv("TEST_DB_MODE") == "ci" {
		// Use CI service database
		dbUser := os.Getenv("POSTGRES_USER")
		dbPass := os.Getenv("POSTGRES_PASSWORD")
		dbName := os.Getenv("POSTGRES_DB")
		dbHost := os.Getenv("POSTGRES_HOST")
		dbPort := os.Getenv("POSTGRES_PORT")

		databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)

		db, err := gorm.Open(postgres.Open(databaseUrl), &gorm.Config{})
		if err != nil {
			return nil, nil, err
		}
		database.DB = db

		// Auto-migrate test database
		err = db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}, &models.Vote{}, &models.SubInvitation{}, &models.Sub{}, &models.SubMembership{}, &models.PasswordResetToken{})
		if err != nil {
			return nil, nil, err
		}

		// No teardown needed for CI
		return db, func() {}, nil
	}
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Printf("Could not construct pool: %s", err)
		return nil, nil, err
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		log.Printf("Could not connect to Docker: %s", err)
		return nil, nil, err
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13",
		Env: []string{
			"POSTGRES_PASSWORD=cortex_pass",
			"POSTGRES_USER=cortex_user",
			"POSTGRES_DB=cortex_db",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Printf("Could not start resource: %s", err)
		return nil, nil, err
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://cortex_user:cortex_pass@%s/cortex_db?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url: ", databaseUrl)

	var db *gorm.DB
	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		db, err = gorm.Open(postgres.Open(databaseUrl), &gorm.Config{})
		if err != nil {
			return err
		}
		database.DB = db
		return db.Exec("SELECT 1").Error
	}); err != nil {
		log.Printf("Could not connect to docker: %s", err)
		return nil, nil, err
	}

	// Auto-migrate test database
	err = db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}, &models.Vote{}, &models.SubInvitation{}, &models.Sub{}, &models.SubMembership{}, &models.PasswordResetToken{})
	if err != nil {
		log.Printf("Failed to migrate test database: %s", err)
		return nil, nil, err
	}

	teardown := func() {
		if err := pool.Purge(resource); err != nil {
			log.Printf("Could not purge resource: %s", err)
		}
	}

	return db, teardown, nil
}

// TestMain is used for integration tests to set up environment
func TestMain(m *testing.M, setupDB func() (*gorm.DB, func())) {
	// Set test environment variables only if not already set
	if os.Getenv("POSTGRES_HOST") == "" {
		os.Setenv("POSTGRES_HOST", "localhost")
	}
	if os.Getenv("POSTGRES_USER") == "" {
		os.Setenv("POSTGRES_USER", "cortex_user")
	}
	if os.Getenv("POSTGRES_PASSWORD") == "" {
		os.Setenv("POSTGRES_PASSWORD", "cortex_pass")
	}
	if os.Getenv("POSTGRES_DB") == "" {
		os.Setenv("POSTGRES_DB", "cortex_db")
	}
	if os.Getenv("POSTGRES_PORT") == "" {
		os.Setenv("POSTGRES_PORT", "5432")
	}

	if setupDB != nil {
		// For integration tests, set up real DB
		db, teardown := setupDB()
		defer teardown()

		database.DB = db

		// Run tests
		code := m.Run()

		os.Exit(code)
	} else {
		// For unit tests, just run
		os.Exit(m.Run())
	}
}

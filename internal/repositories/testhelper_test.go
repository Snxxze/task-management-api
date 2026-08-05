package repositories_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"task-management-api/internal/database"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("task_management_test"),
		postgres.WithUsername("test_user"),
		postgres.WithPassword("test_password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)

	if err != nil {
		t.Skipf("Skipping repository integration test: Docker engine not available or testcontainers error: %v", err)
		return nil, func() {}
	}

	host, err := pgContainer.Host(ctx)
	if err != nil {
		_ = pgContainer.Terminate(ctx)
		t.Fatalf("Failed to get container host: %v", err)
	}

	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		_ = pgContainer.Terminate(ctx)
		t.Fatalf("Failed to get container mapped port: %v", err)
	}

	dsn := fmt.Sprintf("host=%s user=test_user password=test_password dbname=task_management_test port=%s sslmode=disable TimeZone=Asia/Bangkok", host, port.Port())
	db, err := gorm.Open(gormpostgres.Open(dsn), &gorm.Config{})
	if err != nil {
		_ = pgContainer.Terminate(ctx)
		t.Fatalf("Failed to connect to test postgres container: %v", err)
	}

	err = database.AutoMigrate(db)
	if err != nil {
		_ = pgContainer.Terminate(ctx)
		t.Fatalf("Failed to run AutoMigrate on test DB: %v", err)
	}

	teardown := func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Printf("Failed to terminate container: %v", err)
		}
	}

	return db, teardown
}

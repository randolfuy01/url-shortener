package server

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	migrations "url-shortener.com/m/migrations/driver"
)

var pool *pgxpool.Pool
var initOnce sync.Once

// initConnection initializes the shared pgx pool. Safe to call multiple times.
func initConnection() {
	dsn := os.Getenv("DATABASE_URL")
	initOnce.Do(func() {
		cfg, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			panic(fmt.Errorf("failed to parse DATABASE_URL: %w", err))
		}
		p, err := pgxpool.NewWithConfig(context.Background(), cfg)
		if err != nil {
			panic(fmt.Errorf("failed to create pgx pool: %w", err))
		}
		// Verify connectivity early
		if err := p.Ping(context.Background()); err != nil {
			p.Close()
			panic(fmt.Errorf("database ping failed: %w", err))
		}
		pool = p
	})
}

// GetDB returns the initialized pool, initializing on first use.
func GetDB() *pgxpool.Pool {
	if pool == nil {
		initConnection()
	}
	return pool
}

// GetQueries constructs sqlc Queries bound to the shared pool.
func GetQueries() *migrations.Queries {
	return migrations.New(GetDB())
}

// QueriesProviderInterface abstracts sqlc query methods used by handlers for easy mocking in tests.
type QueriesProviderInterface interface {
	GetUserByName(ctx context.Context, name string) (migrations.User, error)
	CreateUser(ctx context.Context, arg migrations.CreateUserParams) (migrations.User, error)
	CreateUrl(ctx context.Context, arg migrations.CreateUrlParams) (migrations.Url, error)
	GetUserURLs(ctx context.Context, userID int64) ([]migrations.Url, error)
}

// ProvideQueries returns a provider for queries. Overridable in tests.
var ProvideQueries = func() QueriesProviderInterface {
	return GetQueries()
}

// CloseDB closes the shared pool. Call from main at shutdown.
func CloseDB() {
	if pool != nil {
		pool.Close()
	}
}

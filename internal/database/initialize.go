package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var Pool *pgxpool.Pool

func Init_DB() {
	err := godotenv.Load(".env")

	// Failed loading environment
	if err != nil {
		log.Fatal("unable to load environment")
	}

	db := os.Getenv("postgres_connection")
	pool, err := pgxpool.New(context.Background(), db)
	if err != nil {
		log.Fatalf("unable to create pool: %v", err)
	}
	Pool = pool
}

// Wrapper function for executing queries
func Execute_Query[T any](item T, operation func(T, *pgxpool.Pool) (bool, error)) (bool, error) {
	result, err := operation(item, Pool)
	if err != nil {
		return false, err
	}
	return result, nil
}

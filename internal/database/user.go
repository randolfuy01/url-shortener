package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID       int
	Name     string
	Password string
}

// Wrapper function for executing queries
func Execute_Query[T any](item T, operation func(T, *pgxpool.Pool) (bool, error)) (bool, error) {
	result, err := operation(item, Pool)
	if err != nil {
		return false, err
	}
	return result, nil
}

// Inserting a new user into the database
func Insert_User(user User, pool *pgxpool.Pool) (bool, error) {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return false, err
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), `INSERT INTO users (name, password) VALUES ($1, $2)`, user.Name, user.Password)
	if err != nil {
		return false, err
	}

	return true, nil
}

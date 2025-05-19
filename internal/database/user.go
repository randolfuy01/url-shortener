package database

import (
	"context"

	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID       int
	Name     string
	Password string
}

// Inserting a new user into the database
func Insert_User(user User, pool *pgxpool.Pool) (bool, error) {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return false, err
	}
	defer conn.Release()

	_, qErr := conn.Exec(context.Background(), `INSERT INTO users (name, password) VALUES ($1, $2)`, user.Name, user.Password)
	if qErr != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return false, errors.New("row already exists")
		}
		return false, err
	}

	return true, nil
}

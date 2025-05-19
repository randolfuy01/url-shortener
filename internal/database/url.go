package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type URL struct {
	ID          int    `db:"id"`
	UserID      int    `db:"user_id"`
	OriginalURL string `db:"original_url"`
	ShortURL    string `db:"short_url"`
}

// Inserting a url into the database
func Create_URL(url URL, pool *pgxpool.Pool) (bool, error) {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return false, err
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), `INSERT INTO urls (user_id, original_url, short_url)`,
		url.UserID, url.OriginalURL, url.ShortURL)

	if err != nil {
		// Unique violation
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return false, errors.New("row already exists")
		}
		return false, err
	}

	return true, nil
}

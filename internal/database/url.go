package database

type URL struct {
	ID          int    `db:"id"`
	UserID      int    `db:"user_id"`
	OriginalURL string `db:"original_url"`
	ShortURL    string `db:"short_url"`
}

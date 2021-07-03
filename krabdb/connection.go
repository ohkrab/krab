package krabdb

import (
	"os"

	"github.com/jmoiron/sqlx"
)

// WithConnection connects to database and performs actions.
// After that connection is closed.
func WithConnection(f func(db *sqlx.DB) error) error {
	db, err := sqlx.Connect("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer db.Close()

	return f(db)
}

package krabdb

import (
	"github.com/jmoiron/sqlx"
	"github.com/ohkrab/krab/krabenv"
)

type Connection interface {
	Get(f func(db DB) error) error
}

type DefaultConnection struct{}

func (d *DefaultConnection) Get(f func(db DB) error) error {
	db, err := Connect(krabenv.DatabaseURL())
	if err != nil {
		return err
	}
	defer db.Close()

	return f(&Instance{database: db})
}

func Connect(connectionString string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", connectionString)
	if err != nil {
		return nil, err
	}
	return db, nil
}

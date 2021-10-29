package krabdb

import (
	"github.com/ohkrab/krab/krabenv"
)

type Connection interface {
	Get(f func(db *DB) error) error
}

type DefaultConnection struct{}

func (d *DefaultConnection) Get(f func(db *DB) error) error {
	db, err := Connect(krabenv.DatabaseURL())
	if err != nil {
		return err
	}
	defer db.Close()

	return f(db)
}

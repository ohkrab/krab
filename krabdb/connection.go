package krabdb

import (
	"net/url"

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

type SwitchableDatabaseConnection struct {
	DatabaseName string
}

func (d *SwitchableDatabaseConnection) Get(f func(db DB) error) error {
	purl, err := url.Parse(krabenv.DatabaseURL())
	if err != nil {
		return err
	}
	if purl.Path != "" {
		purl.Path = d.DatabaseName
	}
	db, err := Connect(purl.String())
	if err != nil {
		return err
	}
	defer db.Close()

	return f(&Instance{database: db})
}

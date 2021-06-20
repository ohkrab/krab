package krab

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/afero"
)

// mockParser expects args: "path", "content", "path2", "content2", ...
func mockParser(pathContentPair ...string) *Parser {
	memfs := afero.NewMemMapFs()

	for i := 1; i < len(pathContentPair); i += 2 {
		path := pathContentPair[i-1]
		content := pathContentPair[i]
		afero.WriteFile(
			memfs,
			path,
			[]byte(content),
			0644,
		)
	}

	p := NewParser()
	p.fs = afero.Afero{Fs: memfs}
	return p
}

func withPg(t *testing.T, f func(*pgx.Conn)) {
	pool, err := pgxpool.Connect(
		context.Background(),
		"postgres://krab:secret@localhost:5432/krab?sslmode=disable",
	)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		t.Fatalf("Failed to acquire conn: %v", err)
	}
	defer conn.Release()

	f(conn.Conn())
}

func withPgSqlx(t *testing.T, f func(db *sqlx.DB)) {
	db, err := sqlx.Connect(
		"pgx",
		"postgres://krab:secret@localhost:5432/krab?sslmode=disable",
	)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Fatalf("Failed to ping db: %v", err)
	}

	f(db)
}

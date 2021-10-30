package spec

import (
	"context"
	"database/sql"
	"io"
	"os"
	"strings"
	"testing"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/ohkrab/krab/krabdb"
)

type testDB struct {
	db       *sqlx.DB
	recorder io.StringWriter
}

func (d *testDB) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	d.recorder.WriteString(query)
	return d.db.SelectContext(ctx, dest, query, args...)
}

func (d *testDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	d.recorder.WriteString(query)
	return d.db.QueryxContext(ctx, query, args...)
}

func (d *testDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	d.recorder.WriteString(query)
	return d.db.ExecContext(ctx, query, args...)
}

func withPg(t *testing.T, f func(db *testDB)) {
	db, err := sqlx.Connect("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	defer cleanDb(db)
	f(&testDB{db: db, recorder: &strings.Builder{}})
}

func cleanDb(db krabdb.ExecerContext) {
	_, err := db.ExecContext(context.TODO(), `
DO 
$$ 
  DECLARE 
    r RECORD;
BEGIN
  FOR r IN 
    (
      SELECT table_schema, table_name 
        FROM information_schema.tables 
       WHERE table_schema NOT IN ('pg_catalog', 'information_schema')
    ) 
  LOOP
     EXECUTE 'DROP TABLE ' || quote_ident(r.table_schema) || '.' || quote_ident(r.table_name) || ' CASCADE';
  END LOOP;
END
$$`)

	if err != nil {
		panic(err)
	}
}

func sqlxRowsMapScan(rows *sqlx.Rows) []map[string]interface{} {
	res := []map[string]interface{}{}
	for rows.Next() {
		row := map[string]interface{}{}
		rows.MapScan(row)
		res = append(res, row)
	}

	return res
}

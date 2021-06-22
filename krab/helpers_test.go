package krab

import (
	"os"
	"testing"

	_ "github.com/jackc/pgx/v4"
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

func withPg(t *testing.T, f func(db *sqlx.DB)) {
	db, err := sqlx.Connect("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Fatalf("Failed to ping db: %v", err)
	}
	defer cleanDb(db)

	f(db)
}

func cleanDb(db *sqlx.DB) {
	db.MustExec(`
DO 
$$ 
  DECLARE 
    r RECORD;
BEGIN
  FOR r IN 
    (
      SELECT table_name 
        FROM information_schema.tables 
       WHERE table_schema = 'public'
    ) 
  LOOP
     EXECUTE 'DROP TABLE ' || quote_ident(r.table_name) || ' CASCADE';
  END LOOP;
END
$$`)
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

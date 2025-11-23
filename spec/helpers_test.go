package spec

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"strings"
	"testing"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ohkrab/krab/ferro"
	"github.com/ohkrab/krab/fmtx"
	"github.com/ohkrab/krab/krab"
	"github.com/ohkrab/krab/krabdb"
	"github.com/ohkrab/krab/plugins"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/wzshiming/ctc"
)

type cliMock struct {
	connection *mockDBConnection
	// config        *krab.Config
	app      *ferro.App
	exitCode int
	fs       afero.Afero
	id       string
	stdout   *bytes.Buffer
	stderr   *bytes.Buffer
	T        *testing.T
}

func NewTestCLI(t *testing.T) *cliMock {
	memfs := afero.NewMemMapFs()

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	app := &ferro.App{
		Logger: fmtx.New(stdout, stderr),
	}

	return &cliMock{
		T:      t,
		id:     uuid.Must(uuid.NewV7()).String(),
		fs:     afero.Afero{Fs: memfs},
		app:    app,
		stdout: stdout,
		stderr: stderr,
	}
}

func (c *cliMock) DefaultDatabase() {
	c.Files(
		"config.fyml",
		`
apiVersion: drivers/v1
kind: Driver
metadata:
  name: test
spec:
  driver: postgresql
  config:
    dsn: postgres://test:test@localhost:5433/test
        `,
	)
}

func (c *cliMock) Files(pathContentPair ...string) {
	for i := 1; i < len(pathContentPair); i += 2 {
		path := pathContentPair[i-1]
		content := pathContentPair[i]
		afero.WriteFile(
			c.fs,
			path,
			[]byte(content),
			0644,
		)
	}
}

func (m *cliMock) AssertRun(args ...string) bool {
	m.setup(args)
	m.exitCode = m.app.Run(args)

	if assert.Equal(m.T, 0, m.exitCode, "Exit code should be eql to 0") {
		return true
	} else {
		fmt.Println("statements debug:")
		for _, sql := range m.connection.recorder {
			fmt.Println("---")
			fmt.Println(ctc.ForegroundBrightRed, sql, ctc.Reset)
		}
		fmt.Println("---")
		fmt.Println(ctc.ForegroundRed, m.stderr.String(), ctc.Reset)
		fmt.Println(ctc.ForegroundRed, m.stdout.String(), ctc.Reset)
		return false
	}
}

func (m *cliMock) setup(args []string) {
	// templates := tpls.New(template.FuncMap{})
	registry := plugins.New()
	registry.RegisterAll()
	// filesystem := config.NewFilesystem("/")
	// runner := run.New(filesystem, templates, registry, m.app.Logger)

	m.connection = &mockDBConnection{
		recorder:         []string{},
		assertedSQLIndex: 0,
	}
	memfs := afero.NewMemMapFs()
	m.fs = afero.Afero{Fs: memfs}
}

func (m *cliMock) Teardown() {
	err := m.connection.Get(func(db krabdb.DB) error {
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

		return err
	})

	if err != nil {
		panic(err)
	}
}

func (m *cliMock) RefuteRun(t *testing.T, args []string) bool {
	m.setup(args)
	m.exitCode = m.app.Run(args)

	return assert.Equal(t, 1, m.exitCode, "Exit code should be greather than 0")
}

func (m *cliMock) AssertOutputContains(t *testing.T, output string) bool {
	return assert.Contains(
		t,
		strings.TrimSpace(m.stdout.String()),
		strings.TrimSpace(output),
		"Output mismatch",
	)
}

func (m *cliMock) AssertUiErrorOutputContains(t *testing.T, output string) bool {
	return assert.Contains(
		t,
		strings.TrimSpace(m.stderr.String()),
		strings.TrimSpace(output),
		"UI error output mismatch",
	)
}

func (m *cliMock) AssertSchemaMigrationTableMissing(t *testing.T, schema string) bool {
	err := m.connection.Get(func(db krabdb.DB) error {
		_, err := krab.NewSchemaMigrationTable(schema).SelectAll(context.TODO(), db)
		return err
	})
	if assert.Error(t, err, "AssertSchemaMigrationTableMissing expects error") {
		return assert.Contains(
			t,
			fmt.Sprintf(`relation "%s.schema_migrations" does not exist`, schema),
			err.Error(),
		)
	}

	return false
}

func (m *cliMock) AssertSchemaMigrationTable(t *testing.T, schema string, expectedVersions ...string) bool {
	var versions []krab.SchemaMigration

	err := m.connection.Get(func(db krabdb.DB) error {
		vers, err := krab.NewSchemaMigrationTable(schema).SelectAll(context.TODO(), db)
		versions = vers
		return err
	})
	if assert.NoError(t, err) {
		if assert.Equal(t, len(versions), len(expectedVersions), "Scheme versions count mismatch") {
			for i, v := range expectedVersions {
				if !assert.Equal(t, versions[i].Version, v) {
					return false
				}
			}

			return true
		} else {
			return false
		}
	}

	return false
}

// AssertSQLContains compares expected query with all recoreded queries.
// Assertions must happen in order the queries are executed, otherwise assertion fails.
func (m *cliMock) AssertSQLContains(t *testing.T, expected string) bool {
	expected = strings.TrimSpace(expected)
	found := -1

	// find matching query and remember last asserted query index
	for i, sql := range m.connection.recorder {
		if strings.Index(sql, expected) != -1 {
			found = i
			if i > m.connection.assertedSQLIndex {
				m.connection.assertedSQLIndex = i
			}
			break
		}
	}

	sql := ""
	if found != -1 {
		sql = m.connection.recorder[found]
		// make sure assertion happen in the correct order
		if found < m.connection.assertedSQLIndex {
			return assert.True(t, false, "Queries asserted in the wrong order")
		}
	}

	return assert.True(t, found != -1, fmt.Sprintf("SQL mismatch:\n%s\nwith:\n%s", expected, sql))
}

func (m *cliMock) FSFiles() map[string][]byte {
	data := map[string][]byte{}
	m.fs.Walk("/", func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			b, err := m.fs.ReadFile(path)
			if err != nil {
				panic(err)
			}
			data[path] = b
		}
		return nil
	})
	return data
}

func (m *cliMock) ResetSQLRecorder() {
	m.connection.assertedSQLIndex = 0
	m.connection.recorder = []string{}
}

func (m *cliMock) Query(t *testing.T, query string) ([]string, []map[string]any) {
	var cols []string
	var vals []map[string]any

	m.connection.Get(func(db krabdb.DB) error {
		rows, err := db.QueryContext(context.TODO(), query)
		defer rows.Close()
		assert.NoError(t, err, fmt.Sprint("Query ", query, " must execute successfully"))

		cols, _ = rows.Columns()
		vals = sqlxRowsMapScan(rows)
		return err
	})

	return cols, vals
}

func (m *cliMock) Insert(t *testing.T, table string, cols string, vals string) bool {
	var err error
	m.connection.Get(func(db krabdb.DB) error {
		_, err = db.ExecContext(
			context.TODO(),
			fmt.Sprintf(
				"INSERT INTO %s(%s) VALUES%s",
				table,
				cols,
				vals,
			),
		)
		return err
	})
	return assert.NoError(t, err, "Insertion must happen")
}

type versionGeneratorMock struct{}

func (g *versionGeneratorMock) Next() string {
	return "20230101"
}

func mockConfig(source string) *krab.Config {
	p := mockParser("src/mock.krab.hcl", source)
	c, err := p.LoadConfigDir("src")
	if err != nil {
		e := fmt.Errorf("Mocking Config failed: %w", err)
		panic(e.Error())
	}
	return c

}

// mockParser expects args: "path", "content", "path2", "content2", ...
func mockParser(pathContentPair ...string) *krab.Parser {
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

	p := krab.NewParser()
	p.FS = afero.Afero{Fs: memfs}
	return p
}

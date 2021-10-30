package spec

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ohkrab/krab/cli"
	"github.com/ohkrab/krab/krab"
	"github.com/ohkrab/krab/krabcli"
	"github.com/ohkrab/krab/krabdb"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

type cliMock struct {
	config        *krab.Config
	app           *krabcli.App
	exitCode      int
	err           error
	uiWriter      bytes.Buffer
	uiErrorWriter bytes.Buffer
	helpWriter    bytes.Buffer
	errorWriter   bytes.Buffer
}

func (m *cliMock) setup(args []string) {
	m.errorWriter = bytes.Buffer{}
	m.helpWriter = bytes.Buffer{}
	m.uiErrorWriter = bytes.Buffer{}
	m.uiWriter = bytes.Buffer{}
	m.app = krabcli.New(cli.New(&m.uiErrorWriter, &m.uiWriter), args, m.config)
	m.app.CLI.ErrorWriter = &m.errorWriter
	m.app.CLI.HelpWriter = &m.helpWriter
}

func (m *cliMock) AssertFailedRun(t *testing.T, args []string) bool {
	m.setup(args)
	m.exitCode, m.err = m.app.Run()

	return assert.Equal(t, 1, m.exitCode, "Exit code should be greather than 0")
}

func (m *cliMock) AssertSuccessfulRun(t *testing.T, args []string) bool {
	m.setup(args)
	m.exitCode, m.err = m.app.Run()

	if assert.NoError(t, m.err, "CLI should run successfully") {
		return assert.Equal(t, 0, m.exitCode, "Exit code should be eql to 0")
	}

	return false
}

func (m *cliMock) AssertOutputContains(t *testing.T, output string) bool {
	return assert.Contains(
		t,
		strings.TrimSpace(m.uiWriter.String()),
		strings.TrimSpace(output),
		"Output mismatch",
	)
}

func (m *cliMock) AssertUiErrorOutputContains(t *testing.T, output string) bool {
	return assert.Contains(
		t,
		strings.TrimSpace(m.uiErrorWriter.String()),
		strings.TrimSpace(output),
		"UI error output mismatch",
	)
}

func (m *cliMock) AssertSchemaMigrationTableMissing(t *testing.T, db krabdb.QueryerContext, schema string) bool {
	_, err := krab.NewSchemaMigrationTable(schema).SelectAll(context.TODO(), db)
	if assert.Error(t, err) {
		return assert.Contains(
			t,
			err.Error(),
			fmt.Sprintf(`relation "%s.schema_migrations" does not exist`, schema),
		)
	}

	return false
}

func (m *cliMock) AssertSchemaMigrationTable(t *testing.T, db krabdb.QueryerContext, schema string, expectedVersions ...string) bool {
	versions, err := krab.NewSchemaMigrationTable(schema).SelectAll(context.TODO(), db)
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

func (m *cliMock) Query(t *testing.T, db krabdb.QueryerContext, query string) ([]string, []map[string]interface{}) {
	rows, err := db.QueryContext(context.TODO(), query)
	assert.NoError(t, err, fmt.Sprint("Query ", query, " must execute successfully"))
	defer rows.Close()

	cols, _ := rows.Columns()
	vals := sqlxRowsMapScan(rows)

	return cols, vals
}

func (m *cliMock) Insert(t *testing.T, db krabdb.ExecerContext, table string, cols string, vals string) bool {
	_, err := db.ExecContext(
		context.TODO(),
		fmt.Sprintf(
			"INSERT INTO %s(%s) VALUES%s",
			table,
			cols,
			vals,
		),
	)
	return assert.NoError(t, err, "Insertion must happen")
}

func mockCli(config *krab.Config) *cliMock {
	mock := &cliMock{
		config: config,
	}

	return mock
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

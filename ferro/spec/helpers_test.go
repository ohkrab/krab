package spec

import (
	"bytes"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ohkrab/krab/ferro"
	"github.com/ohkrab/krab/fmtx"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/wzshiming/ctc"
)

type cliMock struct {
	connection *mockDBConnection
	app        *ferro.App
	exitCode   int
	fs         afero.Afero
	id         string
	stdout     *bytes.Buffer
	stderr     *bytes.Buffer
	T          *testing.T
}

func NewTestCLI(t *testing.T) (*cliMock, func()) {
	dir := t.TempDir()
	osfs := afero.NewOsFs()

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	app := &ferro.App{
		Logger: fmtx.New(stdout, stderr),
		Dir:    dir,
	}

	teardown := func() {
	}

	return &cliMock{
		T:      t,
		id:     uuid.Must(uuid.NewV7()).String(),
		fs:     afero.Afero{Fs: osfs},
		app:    app,
		stdout: stdout,
		stderr: stderr,
	}, teardown
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
		path := filepath.Join(c.app.Dir, pathContentPair[i-1])
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
	m.exitCode = m.app.Run(append([]string{"ferro"}, args...))

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

func (m *cliMock) AssertSchemaMigrationTable(t *testing.T, schema string, expectedVersions ...string) bool {
	panic("AssertSchemaMigrationTable not implemented")
	return true
}

func (m *cliMock) setup(args []string) {
	m.connection = &mockDBConnection{
		recorder:         []string{},
		assertedSQLIndex: 0,
	}
	memfs := afero.NewMemMapFs()
	m.fs = afero.Afero{Fs: memfs}
}

func (m *cliMock) RefuteRun(args ...string) bool {
	m.setup(args)
	m.exitCode = m.app.Run(append([]string{"ferro"}, args...))

	return assert.Greater(m.T, m.exitCode, 0, "Exit code should be greather than 0")
}

func (m *cliMock) AssertOutputContains(t *testing.T, output string) bool {
	val := assert.Contains(
		t,
		strings.TrimSpace(m.stdout.String()),
		strings.TrimSpace(output),
		"Output mismatch",
	)
	if !val {
		t.FailNow()
	}
	return val
}

func (m *cliMock) AssertUiErrorOutputContains(t *testing.T, output string) bool {
	return assert.Contains(
		t,
		strings.TrimSpace(m.stderr.String()),
		strings.TrimSpace(output),
		"UI error output mismatch",
	)
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

type versionGeneratorMock struct{}

func (g *versionGeneratorMock) Next() string {
	return "20230101"
}

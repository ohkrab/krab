package spec

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ohkrab/krab/ferro"
	"github.com/ohkrab/krab/ferro/config"
	"github.com/ohkrab/krab/ferro/plugin"
	"github.com/ohkrab/krab/fmtx"
	"github.com/ohkrab/krab/plugins"
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
		// tempdir automatically cleans after finished test
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

func (c *cliMock) RandomDatabase() func() {
	dbID := fmt.Sprintf("test_%s", strings.ReplaceAll(uuid.NewString(), "-", ""))
	execCtx := plugin.DriverExecutionContext{
		Schema: "public",
	}
	ctx := context.Background()

	driver := plugins.NewPostgreSQLDriver()
	conn, err := driver.Connect(context.Background(), config.DriverConfig{
		"dsn": "postgres://test:test@localhost:5433/test",
	})
	if err != nil {
		c.T.Fatalf("failed to connect to test database: %v", err)
	}
	defer driver.Disconnect(ctx, conn)

	// create database and grant privileges for the test case
	err = conn.Query(execCtx).Exec(ctx, fmt.Sprintf("CREATE DATABASE %s", dbID))
	if err != nil {
		c.T.Fatalf("failed to create database: %v", err)
	}
	err = conn.Query(execCtx).Exec(ctx, fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s TO test;", dbID))
	if err != nil {
		c.T.Fatalf("failed to grant access to database %v", err)
	}

	dbTeardown := func() {
		conn, err := driver.Connect(ctx, config.DriverConfig{
			"dsn": "postgres://test:test@localhost:5433/test",
		})
		if err != nil {
			c.T.Fatalf("failed to connect to test database to perform cleanup: %v", err)
		}
		defer driver.Disconnect(ctx, conn)

		// cleanup
		err = conn.Query(execCtx).Exec(ctx, fmt.Sprintf("DROP DATABASE %s", dbID))
		if err != nil {
			c.T.Fatalf("failed to drop database %v", err)
		}
	}
	c.Files(
		"config.fyml",
		fmt.Sprintf(`
apiVersion: drivers/v1
kind: Driver
metadata:
  name: test
spec:
  driver: postgresql
  config:
    dsn: postgres://test:test@localhost:5433/%s
        `, dbID),
	)

	return dbTeardown
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
}

func (m *cliMock) RefuteRun(args ...string) bool {
	m.setup(args)
	m.exitCode = m.app.Run(append([]string{"ferro"}, args...))

	return assert.Greater(m.T, m.exitCode, 0, "Exit code should be greather than 0")
}

func (m *cliMock) AssertOutputContains(t *testing.T, output string) bool {
	s := fmtx.StripANSI(m.stdout.String())
	s = fmtx.Squish(s)
	val := assert.Contains(
		t,
		strings.TrimSpace(s),
		strings.TrimSpace(output),
		"Output mismatch",
	)
	if !val {
		t.Fatalf("Captured:\n%s", m.stdout.String())
	}
	return val
}

func (m *cliMock) AssertOutputNotContains(t *testing.T, output string) bool {
	val := assert.NotContains(
		t,
		strings.TrimSpace(fmtx.StripANSI(m.stdout.String())),
		strings.TrimSpace(output),
		"Output mismatch",
	)
	if !val {
		t.Fatalf("Captured:\n%s", m.stdout.String())
	}
	return val
}

func (m *cliMock) ResetAllOutputs() {
	m.stdout.Reset()
	m.stderr.Reset()
}

func (m *cliMock) ResetDriverOutputs() {
	m.connection.assertedSQLIndex = 0
	m.connection.recorder = []string{}
}

type versionGeneratorMock struct{}

func (g *versionGeneratorMock) Next() string {
	return "20230101"
}

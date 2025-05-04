package run

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/ohkrab/krab/ferro/config"
	"github.com/ohkrab/krab/plugins"
)

type TestDB struct {
	driver             *plugins.PostgreSQLDriver
	masterConn         *plugins.PostgreSQLDriverConnection
	clientDriverConfig config.DriverConfig
	clientDatabaseName string
	clientDriverName   string
	clear              func()
	fymlFileName       string
	fymlFileContent    string
}

// createTestDB is using master connection to create test databases for each test case.
// it also provides FYML file definitions to avoid repeation when defining ferro files.
func createTestDB(t *testing.T, ctx context.Context) *TestDB {
	drv := plugins.PostgreSQLDriver{}
	conn, err := drv.Connect(ctx, config.DriverConfig{"dsn": "postgres://test:test@localhost:5433/test"})
	if err != nil {
		t.Error(err)
	}
	id := uuid.Must(uuid.NewV7()).String()
	dbname := fmt.Sprintf("ferro_test_%s", strings.ReplaceAll(id, "-", "_"))
	db := TestDB{
		driver:             &drv,
		masterConn:         conn.(*plugins.PostgreSQLDriverConnection),
		clientDatabaseName: dbname,
		clientDriverConfig: config.DriverConfig{"dsn": fmt.Sprintf("postgres://test:test@localhost:5433/%s", dbname)},
		clientDriverName:   "ferro_test",
		fymlFileName:       ".ferro/test_driver.fyml",
		fymlFileContent: fmt.Sprintf(`
apiVersion: drivers/v1
kind: Driver
metadata:
  name: ferro_test
spec:
  driver: postgresql
  config:
    dsn: postgres://
    user: test
    password: test
    db: %s
    port: 5433
`, dbname),
	}

	_, err = db.masterConn.Conn.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s", dbname))
	if err != nil {
		t.Error(fmt.Errorf("failed to create test db: %w", err))
	}
	db.clear = func() {
		_, err := db.masterConn.Conn.Exec(ctx, fmt.Sprintf("DROP DATABASE %s", dbname))
		if err != nil {
			t.Error(fmt.Errorf("failed to drop test db: %w", err))
		}
		err = drv.Disconnect(ctx, db.masterConn)
		if err != nil {
			t.Error(err)
		}

	}
	return &db
}

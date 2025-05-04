package run

import (
	"context"
	"testing"
	"text/template"

	"github.com/ohkrab/krab/ferro/config"
	"github.com/ohkrab/krab/plugins"
	"github.com/ohkrab/krab/tpls"
	"github.com/qbart/expecto/expecto"
)

func TestRunner_MigrationAuditLog(t *testing.T) {
	db := createTestDB(t, context.Background())
	defer db.clear()
	_, dir, fsCleanup := expecto.TempFS(
		db.fymlFileName,
		db.fymlFileContent,
	)
	defer fsCleanup()

	fs := config.NewFilesystem(dir)
	registry := plugins.New()
	registry.RegisterAll()
    templates := tpls.New(template.FuncMap{})
	runner := New(fs, templates, registry)

	cmd := &CommandMigrateUp{
		Driver: db.clientDriverName,
		Set:    "public",
	}

    err := runner.Execute(context.Background(), cmd)
    expecto.NoErr(t, "runner execution", err)

	// builder := NewBuilder(fs, parsed, plugins.New())
	// cfg, errs := builder.BuildConfig()
	// expecto.NotNil(t, "build errors", errs)
	// expecto.Eq(t, "number of errors", len(errs.Errors), 0)
}

func TestRunner_MigrationAuditEntries(t *testing.T) {
	ctx := context.Background()
	db := createTestDB(t, ctx)
	defer db.clear()

	// Create a temporary filesystem with a migration file
	migrationContent := `
kind: Migration
apiVersion: ferro/v1
metadata:
  name: create_test_table
spec:
  version: "20250504120000"
  run:
    up:
      sql: "CREATE TABLE test_table (id SERIAL PRIMARY KEY, name TEXT)"
    down:
      sql: "DROP TABLE test_table"
---
kind: MigrationSet
apiVersion: ferro/v1
metadata:
  name: public
spec:
  namespace:
    schema: public
    prefix: ""
  migrations:
    - create_test_table
`
	_, dir, fsCleanup := expecto.TempFS(
		"migrations.yaml", migrationContent,
	)
	defer fsCleanup()

	// Setup runner
	fs := config.NewFilesystem(dir)
	registry := plugins.New()
	registry.RegisterAll()
	templates := tpls.New(template.FuncMap{})
	runner := New(fs, templates, registry)

	// Execute migration up command
	cmd := &CommandMigrateUp{
		Driver: db.clientDriverName,
		Set:    "public",
	}
	err := runner.Execute(ctx, cmd)
	expecto.NoErr(t, "runner execution", err)

	// Get driver connection to check audit log entries
	driverInstance, err := runner.getDriverInstance(ctx, db.clientDriverName)
	expecto.NoErr(t, "getting driver instance", err)
	
	conn, err := driverInstance.Driver.Connect(ctx, driverInstance.Config.Spec.Config)
	expecto.NoErr(t, "connecting to database", err)
	defer driverInstance.Driver.Disconnect(ctx, conn)

	// Check audit log entries
	execCtx := plugin.DriverExecutionContext{
		Schema: "public",
		Prefix: "",
	}
	
	logs, err := conn.ReadAuditLogs(ctx, execCtx)
	expecto.NoErr(t, "reading audit logs", err)
	
	// Verify we have at least one audit log entry
	expecto.True(t, "audit log has entries", len(logs) > 0)
	
	// Verify the migration version is in the audit log
	found := false
	for _, log := range logs {
		if data, ok := log.Data["version"].(string); ok && data == "20250504120000" {
			found = true
			break
		}
	}
	expecto.True(t, "migration version found in audit log", found)
	
	// Verify the test_table was created
	var tableExists bool
	err = db.driver.(*plugins.PostgreSQLDriver).Conn.QueryRow(ctx, 
		"SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'test_table')").Scan(&tableExists)
	expecto.NoErr(t, "checking if table exists", err)
	expecto.True(t, "test_table exists", tableExists)
}

// TODO: replicate this:

// func TestActionMigrateUp(t *testing.T) {
// 	c := mockCli(mockConfig(`
// migration "do_nothing" {
//   version = "v1"
//
//   up {}
//   down {}
// }
//
// migration_set "public" {
//   migrations = [migration.do_nothing]
// }
// `))
// 	defer c.Teardown()
// 	c.AssertSuccessfulRun(t, []string{"migrate", "up", "public"})
// 	c.AssertOutputContains(t, "\x1b[0;32mOK  \x1b[0mv1 do_nothing")
// 	c.AssertSchemaMigrationTable(t, "public", "v1")
// }
//
// func TestActionMigrateUpWithError(t *testing.T) {
// 	c := mockCli(mockConfig(`migration_set "public" { migrations = [] }`))
// 	defer c.Teardown()
// 	c.AssertSuccessfulRun(t, []string{"migrate", "up", "public"})
//
// 	c = mockCli(mockConfig(`
// migration "do_nothing" {
//   version = "v1"
//
//   up { sql = "SELECT invalid" }
//   down {}
// }
//
// migration_set "public" {
//   migrations = [migration.do_nothing]
// }
// `))
//
// 	c.AssertFailedRun(t, []string{"migrate", "up", "public"})
// 	c.AssertOutputContains(t, "\x1b[0;31mERR \x1b[0mv1 do_nothing")
// 	c.AssertUiErrorOutputContains(t,
// 		`column "invalid" does not exist`,
// 	)
// 	c.AssertSchemaMigrationTable(t, "public")
// }

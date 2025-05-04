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

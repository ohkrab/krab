package spec

import (
	"testing"
)

func TestActionMigrateUp(t *testing.T) {
	cli, teardown := NewTestCLI(t)
	defer teardown()

	dbTeardown := cli.RandomDatabase()
	defer dbTeardown()

	cli.Files(
		"set.fyml",
		`
apiVersion: migrations/v1
kind: MigrationSet
metadata:
  name: public
spec:
  namespace:
    name: public
  migrations:
    - create_animals
        `,
		"01_create_animals.fyml",
		`
apiVersion: migrations/v1
kind: Migration
metadata:
  name: create_animals
spec:
  version: "v1"
  run:
    up:
      sql: CREATE TABLE animals();
    down:
      sql: DROP TABLE animals;
`,
	)

	cli.AssertRun("migrate", "status", "--set", "public", "--driver", "test")
	cli.AssertOutputContains(t, "pending v1 create_animals")
	cli.ResetAllOutputs()

	cli.AssertRun("migrate", "up", "--set", "public", "--driver", "test")
	cli.AssertOutputNotContains(t, "No pending migrations")
	cli.AssertOutputContains(t, "Applied successfully")
	cli.ResetAllOutputs()

	cli.AssertRun("migrate", "up", "--set", "public", "--driver", "test")
	cli.AssertOutputContains(t, "No pending migrations")
	cli.ResetAllOutputs()

	cli.AssertRun("migrate", "status", "--set", "public", "--driver", "test")
	cli.AssertOutputContains(t, "completed v1 create_animals")
	cli.ResetAllOutputs()

	audit := cli.Audit("test", "public")
	audit.AssertCount(2)
	audit.Assert(0, auditLog{
		ID:    1,
		Event: "migration.up.started",
		Data: map[string]any{
			"migration": "create_animals",
			"set":       "public",
			"version":   "v1",
		},
		Metadata: map[string]any{},
	})
	audit.Assert(1, auditLog{
		ID:    2,
		Event: "migration.up.completed",
		Data: map[string]any{
			"migration": "create_animals",
			"set":       "public",
			"version":   "v1",
		},
		Metadata: map[string]any{},
	})
	// TODO:
	// - check if table exists
	// - check if table does not exists
	// - fix english with asserts Contains? NotContain? be consistent

	// cli.AssertRecordsCount("test", "public", 0)
	// ...
	// cli.AssertRecordsCount("test", "public", 1)

	// audit := cli.UseAudit("test", "public")
	// ...
	// audit.AssertEvent(1, "pending")

	// query := cli.UseQuery("test", "public")
	// query.AssertTableExists("animals")
	// query.AssertCount("table", 1)
}

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

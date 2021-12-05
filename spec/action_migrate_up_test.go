package spec

import (
	"testing"
)

func TestActionMigrateUp(t *testing.T) {
	c := mockCli(mockConfig(`
migration "do_nothing" {
  version = "v1"

  up {}
  down {}
}

migration_set "public" {
  migrations = [migration.do_nothing]
}
`))
	defer c.Teardown()
	c.AssertSuccessfulRun(t, []string{"migrate", "up", "public"})
	c.AssertOutputContains(t, "\x1b[0;32mOK  \x1b[0mv1 do_nothing")
	c.AssertSchemaMigrationTable(t, "public", "v1")
}

func TestActionMigrateUpWithError(t *testing.T) {
	c := mockCli(mockConfig(`migration_set "public" { migrations = [] }`))
	defer c.Teardown()
	c.AssertSuccessfulRun(t, []string{"migrate", "up", "public"})

	c = mockCli(mockConfig(`
migration "do_nothing" {
  version = "v1"

  up { sql = "SELECT invalid" }
  down {}
}

migration_set "public" {
  migrations = [migration.do_nothing]
}
`))

	c.AssertFailedRun(t, []string{"migrate", "up", "public"})
	c.AssertOutputContains(t, "\x1b[0;31mERR \x1b[0mv1 do_nothing")
	c.AssertUiErrorOutputContains(t,
		`column "invalid" does not exist`,
	)
	c.AssertSchemaMigrationTable(t, "public")
}

package spec

import (
	"testing"
)

func TestActionMigrateUpHooks(t *testing.T) {
	c := mockCli(mockConfig(`
migration "create_animals" {
  version = "v1"

  up   { sql = "CREATE TABLE animals(name VARCHAR)" }
  down { sql = "DROP TABLE animals" }
}

migration_set "public" {
  schema = "tenants"

  migrations = [migration.create_animals]
}
`))
	defer c.Teardown()

	c.AssertSuccessfulRun(t, []string{"migrate", "up", "public"})
	c.AssertSchemaMigrationTableMissing(t, "public")
	c.AssertSchemaMigrationTable(t, "tenants", "v1")
}

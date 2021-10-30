package spec

import (
	"testing"
)

func TestActionMigrateDownHooks(t *testing.T) {
	c := mockCli(mockConfig(`
migration "create_animals" {
  version = "v1"

  up   { sql = "CREATE TABLE animals(name VARCHAR)" }
  down { sql = "DROP TABLE animals" }
}

migration "add_column" {
  version = "v2"

  up   { sql = "ALTER TABLE animals ADD COLUMN emoji VARCHAR" }
  down { sql = "ALTER TABLE animals DROP COLUMN emoji" }
}

migration_set "tenants" {
  schema = "tenants"

  migrations = [migration.create_animals, migration.add_column]
}
`))
	defer c.Teardown()

	c.AssertSuccessfulRun(t, []string{"migrate", "up", "tenants"})
	c.AssertSchemaMigrationTableMissing(t, "public")
	c.AssertSchemaMigrationTable(t, "tenants", "v1", "v2")

	c.AssertSuccessfulRun(t, []string{"migrate", "down", "tenants", "-version", "v2"})
	c.AssertSchemaMigrationTable(t, "tenants", "v1")
}

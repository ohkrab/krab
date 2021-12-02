package spec

import (
	"testing"
)

func TestActionMigrateUpArguments(t *testing.T) {
	c := mockCli(mockConfig(`
migration "create_animals" {
  version = "v1"

  up   { sql = "CREATE TABLE animals(name VARCHAR)" }
  down { sql = "DROP TABLE animals" }
}

migration_set "public" {
  arguments {
    arg "schema" {}
  }

  schema = "{{.Args.schema}}"

  migrations = [migration.create_animals]
}
`))
	defer c.Teardown()

	c.AssertSuccessfulRun(t, []string{"migrate", "up", "public", "-schema", "custom"})
	c.AssertSchemaMigrationTableMissing(t, "public")
	c.AssertSchemaMigrationTable(t, "custom", "v1")
	c.AssertSQLContains(t, `SET search_path TO "custom"`)
}

package spec

import (
	"testing"

	"github.com/ohkrab/krab/krabdb"
)

func TestActionMigrateUpArguments(t *testing.T) {
	withPg(t, func(db *krabdb.DB) {
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

		c.AssertSuccessfulRun(t, []string{"migrate", "up", "public", "-schema", "custom"})
		c.AssertSchemaMigrationTableMissing(t, db, "public")
		c.AssertSchemaMigrationTable(t, db, "custom", "v1")
	})
}

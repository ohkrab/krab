package spec

import (
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestActionMigrateDownArugments(t *testing.T) {
	withPg(t, func(db *sqlx.DB) {
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
  arguments {
    arg "schema" {}
  }

  schema = "{{.Args.schema}}"

  migrations = [migration.create_animals, migration.add_column]
}
`))

		c.AssertSuccessfulRun(t, []string{"migrate", "up", "tenants", "-schema", "tenants"})
		c.AssertSchemaMigrationTableMissing(t, db, "public")
		c.AssertSchemaMigrationTable(t, db, "tenants", "v1", "v2")

		c.AssertSuccessfulRun(t, []string{"migrate", "down", "tenants", "-schema", "tenants", "-version", "v2"})
		c.AssertSchemaMigrationTable(t, db, "tenants", "v1")
	})
}

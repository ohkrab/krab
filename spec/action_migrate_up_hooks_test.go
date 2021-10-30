package spec

import (
	"testing"

	"github.com/ohkrab/krab/krabdb"
)

func TestActionMigrateUpHooks(t *testing.T) {
	withPg(t, func(db *krabdb.DB) {
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

		c.AssertSuccessfulRun(t, []string{"migrate", "up", "public"})
		c.AssertSchemaMigrationTableMissing(t, db, "public")
		c.AssertSchemaMigrationTable(t, db, "tenants", "v1")
	})
}

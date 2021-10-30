package spec

import (
	"testing"
)

func TestActionMigrateDownTransactions(t *testing.T) {
	withPg(t, func(db *testDB) {
		c := mockCli(mockConfig(`
migration "create_animals" {
  version = "v1"

  up   { sql = "CREATE TABLE animals(name VARCHAR)" }
  down { sql = "DROP TABLE animals" }
}

migration "add_index" {
  version = "v2"

  transaction = false

  up   { sql = "CREATE INDEX CONCURRENTLY idx ON animals(name)" }
  down { sql = "DROP INDEX idx" }
}

migration_set "public" {
  migrations = [migration.create_animals, migration.add_index]
}
`))

		c.AssertSuccessfulRun(t, []string{"migrate", "up", "public"})
		c.AssertSchemaMigrationTable(t, db, "public", "v1", "v2")

		c.AssertSuccessfulRun(t, []string{"migrate", "down", "public", "-version", "v2"})
		c.AssertSchemaMigrationTable(t, db, "public", "v1")
	})
}

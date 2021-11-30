package spec

import (
	"testing"
)

func TestActionMigrateDslStatementOrdering(t *testing.T) {
	c := mockCli(mockConfig(`
migration "create_animals" {
  version = "v1"

  up {
	create_table "animals" {
	  column "id" "bigint" {}
	}

	create_index "animals" "idx_id" {
	  columns = ["id"]
	}

	sql = "ALTER INDEX idx_id RENAME TO idx_new"

	drop_index "idx_new" {}

	drop_table "animals" {}
  }

  down {}
}

migration_set "animals" {
  migrations = [
    migration.create_animals
  ]
}
`))
	defer c.Teardown()
	if c.AssertSuccessfulRun(t, []string{"migrate", "up", "animals"}) {
		c.AssertSchemaMigrationTable(t, "public", "v1")
		c.AssertSQLContains(t, `
CREATE TABLE "animals"(
  "id" bigint
)
	`)
		c.AssertSQLContains(t, `
CREATE INDEX "idx_id" ON "animals"
	`)
		c.AssertSQLContains(t, `
ALTER INDEX idx_id RENAME TO idx_new
	`)
		c.AssertSQLContains(t, `
DROP INDEX "idx_new"
	`)
		c.AssertSQLContains(t, `
DROP TABLE "animals"
	`)
	}
}

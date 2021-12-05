package spec

import (
	"testing"
)

func TestActionMigrateDslIndex(t *testing.T) {
	c := mockCli(mockConfig(`
migration "create_animals" {
  version = "v1"
  transaction = false

  up {
	create_table "animals" {
	  column "id" "bigint" {}

	  column "name" "varchar" {}
	  
	  column "extinct" "boolean" {}

	  column "weight_kg" "int" {}
	}

	create_index "animals" "idx_uniq_name" {
	  unique  = true
	  columns = ["name"]
	  using   = "btree"
	  include = ["weight_kg"]
	}

	create_index "animals" "idx_heavy_animals" {
	  columns      = ["weight_kg"]
	  where        = "weight_kg > 5000"
	  concurrently = true
	}
  }

  down {
    drop_index "public.idx_uniq_name" {
	  cascade = true
	}

    drop_index "idx_heavy_animals" {
	  concurrently = true
	}

    drop_table "animals" {}
  }
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
		c.AssertOutputContains(t, "\x1b[0;32mOK  \x1b[0mv1 create_animals")
		c.AssertSQLContains(t, `
CREATE TABLE "animals"(
  "id" bigint,
  "name" varchar,
  "extinct" boolean,
  "weight_kg" int
)
	`)
		c.AssertSQLContains(t, `
CREATE UNIQUE INDEX "idx_uniq_name" ON "animals" USING btree ("name") INCLUDE ("weight_kg")
	`)
		c.AssertSQLContains(t, `
CREATE INDEX CONCURRENTLY "idx_heavy_animals" ON "animals" ("weight_kg") WHERE (weight_kg > 5000)
	`)

		if c.AssertSuccessfulRun(t, []string{"migrate", "down", "animals", "-version", "v1"}) {
			c.AssertSchemaMigrationTable(t, "public")
			c.AssertOutputContains(t, "\x1b[0;32mOK  \x1b[0mv1 create_animals")
			c.AssertSQLContains(t, `DROP INDEX "public"."idx_uniq_name" CASCADE`)
			c.AssertSQLContains(t, `DROP INDEX CONCURRENTLY "idx_heavy_animals"`)
			c.AssertSQLContains(t, `DROP TABLE "animals"`)
		}
	}
}

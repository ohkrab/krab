package spec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActionCustom(t *testing.T) {
	c := mockCli(mockConfig(`
migration "create_animals" {
  version = "v1"

  up   { sql = "CREATE TABLE animals(name VARCHAR)" }
  down { sql = "DROP TABLE animals" }
}

migration "create_animals_view" {
  version = "v2"

  up   { sql = "CREATE MATERIALIZED VIEW anims AS SELECT name FROM animals" }
  down { sql = "DROP MATERIALIZED VIEW anims" }
}

migration "seed_animals" {
  version = "v3"

  up   { sql = "INSERT INTO animals(name) VALUES('Elephant'),('Turtle'),('Cat')" }
  down { sql = "TRUNCATE animals" }
}

migration_set "animals" {
  migrations = [
    migration.create_animals,
	migration.create_animals_view,
	migration.seed_animals,
  ]
}

action "view" "refresh" {
  arguments {
    arg "name" {}
  }

  sql = "REFRESH MATERIALIZED VIEW {{ quote_ident .Args.name }}"
}
`))
	defer c.Teardown()

	c.AssertSuccessfulRun(t, []string{"migrate", "up", "animals"})
	c.AssertSchemaMigrationTable(t, "public", "v1", "v2", "v3")

	_, vals := c.Query(t, "SELECT * FROM anims")
	if assert.Len(t, vals, 0, "No values should be returned") {
		c.AssertSuccessfulRun(t, []string{"action", "view", "refresh", "-name", "anims"})
		_, vals := c.Query(t, "SELECT * FROM anims")
		assert.Len(t, vals, 3, "There should be 3 animals after refresh")
		c.AssertSQLContains(t, `REFRESH MATERIALIZED VIEW "anims"`)
	}

}

package spec

import (
	"testing"
)

func TestActionMigrateStatusArguments(t *testing.T) {
	c := mockCli(mockConfig(`
migration "create_animals" {
  version = "v1"

  up   { sql = "CREATE TABLE animals(name VARCHAR)" }
  down { sql = "DROP TABLE animals" }
}

migration_set "animals" {
  arguments {
    arg "schema" {}
  }

  schema = "{{.Args.schema}}"

  migrations = [migration.create_animals]
}
`))
	defer c.Teardown()

	c.AssertSuccessfulRun(t, []string{"migrate", "up", "animals", "-schema", "custom"})
	c.AssertSuccessfulRun(t, []string{"migrate", "status", "animals", "-schema", "custom"})
	c.AssertOutputContains(t, "\x1b[0;32m+ \x1b[0mv1 create_animals")
}

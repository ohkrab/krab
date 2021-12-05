package spec

import (
	"fmt"
	"testing"

	"github.com/ohkrab/krab/emojis"
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
	c.AssertOutputContains(t, fmt.Sprint(emojis.CheckMark(), " v1 create_animals"))
}

package spec

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/ohkrab/krab/krab"
)

func TestActionMigrateUp(t *testing.T) {
	withPg(t, func(db *sqlx.DB) {
		c := mockCli(mockConfig(`
migration "do_nothing" {
  version = "v1"

  up {}
  down {}
}

migration_set "public" {
  migrations = [migration.do_nothing]
}
`))
		c.AssertSuccessfulRun(t, []string{"migrate", "up", "public"})
		c.AssertOutputContains(t,
			`
do_nothing v1
Done
`,
		)
		c.AssertSchemaMigrationTable(t, db, "public", "v1")
	})
}

func TestActionMigrateUpWithError(t *testing.T) {
	// assert := assert.New(t)

	withPg(t, func(db *sqlx.DB) {
		krab.NewSchemaMigrationTable("public").Init(context.TODO(), db)

		c := mockCli(mockConfig(`
migration "do_nothing" {
  version = "v1"

  up { sql = "SELECT invalid" }
  down {}
}

migration_set "public" {
  migrations = [migration.do_nothing]
}
`))

		c.AssertFailedRun(t, []string{"migrate", "up", "public"})
		c.AssertUiErrorOutputContains(t,
			`column "invalid" does not exist`,
		)
		c.AssertSchemaMigrationTable(t, db, "public")
	})
}

package spec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActionMigrateDown(t *testing.T) {
	assert := assert.New(t)

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

migration_set "public" {
  migrations = [migration.create_animals, migration.add_column]
}
`))
	defer c.Teardown()
	c.AssertSuccessfulRun(t, []string{"migrate", "up", "public"})
	c.AssertOutputContains(t, "\x1b[0;32mOK  \x1b[0mv1 create_animals")
	c.AssertOutputContains(t, "\x1b[0;32mOK  \x1b[0mv2 add_column")
	c.AssertSchemaMigrationTable(t, "public", "v1", "v2")
	c.Insert(t, "animals", "name, emoji", "('Elephant', '🐘')")
	cols, rows := c.Query(t, "SELECT * from animals")

	assert.ElementsMatch([]string{"name", "emoji"}, cols, "Columns must match")
	if assert.Equal(1, len(rows)) {
		assert.Equal("Elephant", rows[0]["name"])
		assert.Equal("🐘", rows[0]["emoji"])
	}

	c.AssertSuccessfulRun(t, []string{"migrate", "down", "public", "-version", "v2"})
	c.AssertOutputContains(t, "\x1b[0;32mOK  \x1b[0mv2 add_column")
	c.AssertSchemaMigrationTable(t, "public", "v1")

	cols, rows = c.Query(t, "SELECT * from animals")

	assert.ElementsMatch([]string{"name"}, cols, "Columns must match")
	if assert.Equal(1, len(rows)) {
		assert.Equal("Elephant", rows[0]["name"])
		assert.Nil(rows[0]["emoji"])
	}
}

func TestActionMigrateDownOnError(t *testing.T) {
	assert := assert.New(t)

	c := mockCli(mockConfig(`
migration "create_animals" {
  version = "v1"

  up   { sql = "CREATE TABLE animals(name VARCHAR)" }
  down { sql = "DROP TABLE animals" }
}

migration "add_column" {
  version = "v2"

  up   { sql = "ALTER TABLE animals ADD COLUMN emoji VARCHAR" }
  down { sql = "ALTER TABLE animals DROP COLUMN emoji; ALTER TABLE animals DROP COLUMN habitat" }
}

migration_set "public" {
  migrations = [migration.create_animals, migration.add_column]
}
`))
	defer c.Teardown()
	c.AssertSuccessfulRun(t, []string{"migrate", "up", "public"})
	c.AssertOutputContains(t, "\x1b[0;32mOK  \x1b[0mv1 create_animals")
	c.AssertOutputContains(t, "\x1b[0;32mOK  \x1b[0mv2 add_column")
	c.AssertSchemaMigrationTable(t, "public", "v1", "v2")
	c.Insert(t, "animals", "name, emoji", "('Elephant', '🐘')")
	cols, rows := c.Query(t, "SELECT * from animals")

	assert.ElementsMatch([]string{"name", "emoji"}, cols, "Columns must match")
	if assert.Equal(1, len(rows)) {
		assert.Equal("Elephant", rows[0]["name"])
		assert.Equal("🐘", rows[0]["emoji"])
	}

	c.AssertFailedRun(t, []string{"migrate", "down", "public", "-version", "v2"})
	c.AssertUiErrorOutputContains(t,
		`column "habitat" of relation "animals" does not exist`,
	)

	// state after
	c.AssertSchemaMigrationTable(t, "public", "v1", "v2")
	cols, rows = c.Query(t, "SELECT * from animals")
	assert.ElementsMatch([]string{"name", "emoji"}, cols, "Columns must match")
}

func TestActionMigrateDownWhenSchemaDoesNotExist(t *testing.T) {
	assert := assert.New(t)

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

migration_set "public" {
  migrations = [migration.create_animals, migration.add_column]
}
`))
	defer c.Teardown()
	c.AssertSuccessfulRun(t, []string{"migrate", "up", "public"})
	c.AssertSchemaMigrationTable(t, "public", "v1", "v2")
	c.AssertSuccessfulRun(t, []string{"migrate", "down", "public", "-version", "v2"})
	c.AssertSchemaMigrationTable(t, "public", "v1")

	c.Insert(t, "animals", "name", "('Crab')")
	_, rows := c.Query(t, "SELECT * from animals")
	if assert.Equal(1, len(rows)) {
		assert.Equal("Crab", rows[0]["name"])
	}

	c.AssertFailedRun(t, []string{"migrate", "down", "public", "-version", "v2"})
	c.AssertSchemaMigrationTable(t, "public", "v1")
	c.AssertUiErrorOutputContains(t,
		`Migration has not been run yet, nothing to rollback`,
	)
}

package spec

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParserWithDuplicatedRefNames(t *testing.T) {
	assert := assert.New(t)

	p := mockParser(
		"src/public.krab.hcl",
		`
migration "abc" {
  version = "2006"
  up { sql = "" }
  down { sql = "" }
}

migration "abc" {
  version = "2006"
  up { sql = "" }
  down { sql = "" }
}
`)
	_, err := p.LoadConfigDir("src")
	if assert.Error(err) {
		assert.Contains(err.Error(), "Migration with the name 'abc' already exists")
	}
}

func TestParserMigrationSetWithDuplicatedRefName(t *testing.T) {
	assert := assert.New(t)

	p := mockParser(
		"src/sets.krab.hcl",
		`
migration_set "abc" {
  migrations = []
}

migration_set "abc" {
  migrations = []
}
`)
	_, err := p.LoadConfigDir("src")
	if assert.Error(err) {
		assert.Contains(err.Error(), "Migration Set with the name 'abc' already exists", "Names must be unique")
	}
}

func TestParserMigrationSetWithMissingMigrationReference(t *testing.T) {
	assert := assert.New(t)

	p := mockParser(
		"src/sets.krab.hcl",
		`
migration_set "abc" {
  migrations = [migration.does_not_exist]
}
`)
	_, err := p.LoadConfigDir("src")
	if assert.Error(err, "Parsing config should fail") {
		assert.Contains(err.Error(), "Migration Set references 'does_not_exist' migration that does not exist", "Missing migration")
	}
}

func TestParserWithMigrationsDefinedInSQLFiles(t *testing.T) {
	assert := assert.New(t)

	p := mockParser(
		"src/migrations.krab.hcl",
		`
migration "abc" {
  version = "2006"
  up {
	sql = file_read("src/up.sql")
  }
  down {
	sql = file_read("src/down.sql")
  }
}
`,
		"src/up.sql",
		"CREATE TABLE abc",
		"src/down.sql",
		"DROP TABLE abc",
	)

	config, err := p.LoadConfigDir("src")
	if assert.NoError(err, "Parsing config should not fail") {

		migration, exists := config.Migrations["abc"]
		if assert.True(exists) {
			assert.Equal(migration.RefName, "abc")
			var up strings.Builder
			var down strings.Builder
			migration.Up.ToSQL(&up)
			migration.Down.ToSQL(&down)
			assert.Equal(up.String(), "CREATE TABLE abc;\n")
			assert.Equal(down.String(), "DROP TABLE abc;\n")
		}
	}
}

func TestParserWithMigrationsDefinedInSQLFilesThatAreMissing(t *testing.T) {
	assert := assert.New(t)

	p := mockParser(
		"src/migrations.krab.hcl",
		`
migration "abc" {
  version = "2006"
  up {
	sql = file_read("src/up.sql")
  }
  down {
	sql = file_read("src/down.sql")
  }
}
`,
	)

	_, err := p.LoadConfigDir("src")
	if assert.Error(err, "Parsing config should fail") {
		assert.Contains(
			err.Error(),
			`Call to function "file_read" failed`,
		)
	}
}

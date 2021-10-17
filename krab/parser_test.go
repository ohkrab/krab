package krab

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	assert := assert.New(t)

	p := mockParser(
		"src/public.krab.hcl",
		`
migration "create_tenants" {
  version = "2006"

  up {
	sql = "CREATE TABLE tenants(name VARCHAR PRIMARY KEY)"
  }

  down {
	sql = "DROP TABLE tenants"
  }
}
`)
	c, err := p.LoadConfigDir("src")
	if assert.NoError(err) {
		migration, _ := c.Migrations["create_tenants"]

		assert.Equal(migration.RefName, "create_tenants")
		assert.Equal(migration.Version, "2006")
		assert.Equal(migration.Up.SQL, "CREATE TABLE tenants(name VARCHAR PRIMARY KEY)")
		assert.Equal(migration.Down.SQL, "DROP TABLE tenants")
	}
}

func TestParserWithoutMigrationDetails(t *testing.T) {
	assert := assert.New(t)

	p := mockParser(
		"src/public.krab.hcl",
		`migration "abc" {
				  version = "2006"
                  up {}
				  down {}
				}`,
	)
	c, err := p.LoadConfigDir("src")
	if assert.NoError(err) {
		migration, _ := c.Migrations["abc"]
		assert.Equal(migration.RefName, "abc")
		assert.Equal(migration.Up.SQL, "")
		assert.Equal(migration.Down.SQL, "")
	}
}

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

func TestParserMigrationSet(t *testing.T) {
	assert := assert.New(t)

	p := mockParser(
		"src/migrations.krab.hcl",
		`
migration "abc" {
  version = "2006"
  up {}
  down {}
}

migration "def" {
  version = "2006"
  up {}
  down {}
}

migration "xyz" {
  version = "2006"
  up {}
  down {}
}
`,
		"src/sets.krab.hcl",
		`
migration_set "public" {
  migrations = [
  	migration.abc,
	migration.def,
  ]
}

migration_set "private" {
  migrations = [migration.xyz]
}
`)
	c, err := p.LoadConfigDir("src")
	assert.NoError(err)

	// public set
	publicSet, _ := c.MigrationSets["public"]
	assert.Equal(publicSet.RefName, "public")
	assert.Equal(len(publicSet.Migrations), 2)
	assert.Equal(publicSet.Migrations[0].RefName, "abc")
	assert.Equal(publicSet.Migrations[1].RefName, "def")

	// private set
	privateSet, _ := c.MigrationSets["private"]
	assert.Equal(privateSet.RefName, "private")
	assert.Equal(len(privateSet.Migrations), 1)
	assert.Equal(privateSet.Migrations[0].RefName, "xyz")

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
			assert.Equal(migration.Up.SQL, "CREATE TABLE abc")
			assert.Equal(migration.Down.SQL, "DROP TABLE abc")
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

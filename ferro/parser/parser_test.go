package parser

import (
	"testing"

	"github.com/qbart/expecto/expecto"
)

func TestParser_SingleMigrationSetWithMigration(t *testing.T) {
	fs, _, cleanup := expecto.TempFS(
		"src/a/b/c/animals.fyml",
		`
apiVersion: migrations/v1
kind: Migration
metadata:
  name: CreateAnimals
spec:
  version: "v1"
  run:
    up:
      sql: "CREATE TABLE animals(name varchar PRIMARY KEY)"
    down:
      sql: "DROP TABLE animals"
---
apiVersion: migrations/v1
kind: MigrationSet
metadata:
  name: public
spec:
  migrations:
    - CreateAnimals
`)
	defer cleanup()

	should := expecto.New(t)

	parsed, err := (&Parser{FS: fs}).LoadConfigDir("src")
	should.NoErr("parsing config", err)

	cfg, err := parsed.BuildConfig()
	should.NoErr("building config", err)

	should.Eq("number of files",
		len(parsed.Files), 1)
	should.Eq("number of chunks",
		len(parsed.Files[0].Chunks), 2)
	should.Eq("number of migrations",
		len(parsed.Files[0].Migrations), 1)
	should.Eq("number of migration sets",
		len(parsed.Files[0].MigrationSets), 1)

	migrations := expecto.Map(t, cfg.Migrations)
	migrations.HasKey("has animals", "CreateAnimals")

	sets := expecto.Map(t, cfg.MigrationSets)
	sets.HasKey("has public", "public")

	// migration
	should.Eq("migration name",
		cfg.Migrations["CreateAnimals"].Metadata.Name,
		"CreateAnimals")
	should.Eq("migration version",
		cfg.Migrations["CreateAnimals"].Spec.Version,
		"v1")
	should.NotNil("migration transaction",
		cfg.Migrations["CreateAnimals"].Spec.Transaction)
	should.True("deafult is true",
		*cfg.Migrations["CreateAnimals"].Spec.Transaction)
	should.Eq("migration up",
		cfg.Migrations["CreateAnimals"].Spec.Run.Up.Sql,
		"CREATE TABLE animals(name varchar PRIMARY KEY)")
	should.Eq("migration down",
		cfg.Migrations["CreateAnimals"].Spec.Run.Down.Sql,
		"DROP TABLE animals")

	// migration set
	should.Eq("migration set name",
		cfg.MigrationSets["public"].Metadata.Name,
		"public")
	should.Eq("migration set migrations",
		cfg.MigrationSets["public"].Spec.Migrations,
		[]string{"CreateAnimals"})
}

func TestParser_WithDuplicatedRefNames(t *testing.T) {
	fs, _, cleanup := expecto.TempFS(
		"src/animals.fyml",
		`
apiVersion: migrations/v1
kind: Migration
metadata:
  name: CreateAnimals
spec:
  version: "v1"
  run:
    up:
      sql: "CREATE TABLE animals(name varchar PRIMARY KEY)"
    down:
      sql: "DROP TABLE animals"
---
apiVersion: migrations/v1
kind: Migration
metadata:
  name: CreateAnimals
spec:
  version: "v2"
  run:
    up:
      sql: "CREATE TABLE habitats(name varchar PRIMARY KEY)"
    down:
      sql: "DROP TABLE habitats"
`)
	defer cleanup()

	should := expecto.New(t)
	parsed, err := (&Parser{FS: fs}).LoadConfigDir("src")

	should.NoErr("parsing config", err)

	_, err = parsed.BuildConfig()
	should.ErrContains("duplicate migration", err, "adding Migration: migration `CreateAnimals` already exists")
}

func TestParser_MigrationSetWithDuplicatedRefName(t *testing.T) {
	fs, _, cleanup := expecto.TempFS(
		"src/sets.fyml",
		`
apiVersion: migrations/v1
kind: MigrationSet
metadata:
  name: public
spec:
  migrations: []
---
apiVersion: migrations/v1
kind: MigrationSet
metadata:
  name: public
spec:
  migrations: []
`)
	defer cleanup()

	should := expecto.New(t)
	parsed, err := (&Parser{FS: fs}).LoadConfigDir("src")

	should.NoErr("parsing config", err)

	_, err = parsed.BuildConfig()
	should.ErrContains("duplicate migration set", err, "adding MigrationSet: migration set `public` already exists")
}

// func TestParserMigrationSetWithMissingMigrationReference(t *testing.T) {
// 	assert := assert.New(t)

// 	parser, cleanup := mockParser(
// 		"src/sets.fyml",
// 		`
// migration_set "abc" {
//   migrations = [migration.does_not_exist]
// }
// `)
// 	defer cleanup()

// 	_, err := parser.LoadConfigDir("src")
// 	if assert.Error(err, "Parsing config should fail") {
// 		assert.Contains(err.Error(), "Migration Set references 'does_not_exist' migration that does not exist", "Missing migration")
// 	}
// }

// func TestParserWithMigrationsDefinedInSQLFiles(t *testing.T) {
// 	assert := assert.New(t)

// 	parser, cleanup := mockParser(
// 		"src/migrations.fyml",
// 		`
// migration "abc" {
//   version = "2006"
//   up {
// 	sql = file_read("src/up.sql")
//   }
//   down {
// 	sql = file_read("src/down.sql")
//   }
// }
// `,
// 		"src/up.sql",
// 		"CREATE TABLE abc",
// 		"src/down.sql",
// 		"DROP TABLE abc",
// 	)
// 	defer cleanup()

// 	config, err := parser.LoadConfigDir("src")
// 	if assert.NoError(err, "Parsing config should not fail") {
// 		if assert.Equal(len(config.Files), 1) {
// 			migration := config.Files[0].Migrations[0]
// 		}
// 		migration := config.Files[0].Migrations[0]
// 		if assert.True(exists) {
// 			assert.Equal(migration.RefName, "abc")
// 			var up strings.Builder
// 			var down strings.Builder
// 			migration.Up.ToSQL(&up)
// 			migration.Down.ToSQL(&down)
// 			assert.Equal(up.String(), "CREATE TABLE abc")
// 			assert.Equal(down.String(), "DROP TABLE abc")
// 		}
// 	}
// }

// func TestParserWithMigrationsDefinedInSQLFilesThatAreMissing(t *testing.T) {
// 	assert := assert.New(t)

// 	parser, cleanup := mockParser(
// 		"src/migrations.fyml",
// 		`
// migration "abc" {
//   version = "2006"
//   up {
// 	sql = file_read("src/up.sql")
//   }
//   down {
// 	sql = file_read("src/down.sql")
//   }
// }
// `,
// 	)
// 	defer cleanup()

// 	_, err := parser.LoadConfigDir("src")
// 	if assert.Error(err, "Parsing config should fail") {
// 		assert.Contains(
// 			err.Error(),
// 			`Call to function "file_read" failed`,
// 		)
// 	}
// }

// func TestParserRecursiveDir(t *testing.T) {
// 	assert := assert.New(t)

// 	parser, cleanup := mockParser(
// 		"src/a.fyml",
// 		`
// migration "abc" {
//   version = "v1"
//   up {}
//   down {}
// }
// `,
// 		"src/nested/b.fyml",
// 		`
// migration "def" {
//   version = "v2"
//   up {}
//   down {}
// }
// `,
// 	)
// 	defer cleanup()

// 	config, err := parser.LoadConfigDir("src")
// 	if assert.NoError(err, "Parsing config should not fail") {
// 		_, abcOk := config.Migrations["abc"]
// 		_, defOk := config.Migrations["def"]

// 		assert.True(abcOk, "`abc` migration exists")
// 		assert.True(defOk, "`def` migration exists")
// 	}

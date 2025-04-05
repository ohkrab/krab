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

	cfg, errs := parsed.BuildConfig()
	should.Nil("build errors", errs)

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

	_, errs := parsed.BuildConfig()
	should.NotNil("build errors", errs)
	should.Eq("number of errors",
		len(errs.Errors), 1)
	should.ErrContains("duplicate migration", errs.Errors[0], "adding Migration: migration `CreateAnimals` already exists")
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

	_, errs := parsed.BuildConfig()
	should.NotNil("build errors", errs)
	should.Eq("number of errors",
		len(errs.Errors), 1)
	should.ErrContains("duplicate migration set", errs.Errors[0], "adding MigrationSet: migration set `public` already exists")
}

func TestParser_MigrationSetWithMissingMigrationReference(t *testing.T) {
	should := expecto.New(t)

	fs, _, cleanup := expecto.TempFS(
		"src/sets.fyml",
		`
apiVersion: migrations/v1
kind: MigrationSet
metadata:
  name: public
spec:
  migrations: ["DoesNotExist"]
`)
	defer cleanup()

	parsed, err := (&Parser{FS: fs}).LoadConfigDir("src")
	should.NoErr("parsing config", err)

	_, errs := parsed.BuildConfig()
	should.NotNil("build errors", errs)
	should.Eq("number of errors",
		len(errs.Errors), 1)
	should.ErrContains("missing migration", errs.Errors[0],
		"invalid reference: Migration `DoesNotExist` (referenced by MigrationSet `public`) does not exist")
}

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

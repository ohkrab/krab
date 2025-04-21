package run

import (
	"testing"

	"github.com/ohkrab/krab/ferro/config"
	"github.com/ohkrab/krab/plugins"
	"github.com/qbart/expecto/expecto"
)

func TestParser_SingleMigrationSetWithMigration(t *testing.T) {
	_, dir, cleanup := expecto.TempFS(
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

	fs := config.NewFilesystem(dir)
	parser := config.NewParser(fs)
	parsed, err := parser.LoadAndParse()
	expecto.NoErr(t, "parsing config", err)

	builder := NewBuilder(fs, parsed, plugins.New())
	cfg, errs := builder.BuildConfig()
	expecto.NotNil(t, "build errors", errs)
	expecto.Eq(t, "number of errors", len(errs.Errors), 0)

	expecto.Eq(t, "number of files",
		len(parsed.Files), 1)
	expecto.Eq(t, "number of chunks",
		len(parsed.Files[0].Chunks), 2)
	expecto.Eq(t, "number of migrations",
		len(parsed.Files[0].Migrations), 1)
	expecto.Eq(t, "number of migration sets",
		len(parsed.Files[0].MigrationSets), 1)

	migrations := expecto.Map(t, cfg.Migrations)
	migrations.HasKey("has animals", "CreateAnimals")

	sets := expecto.Map(t, cfg.MigrationSets)
	sets.HasKey("has public", "public")

	// migration
	expecto.Eq(t, "migration name",
		cfg.Migrations["CreateAnimals"].Metadata.Name,
		"CreateAnimals")
	expecto.Eq(t, "migration version",
		cfg.Migrations["CreateAnimals"].Spec.Version,
		"v1")
	expecto.NotNil(t, "migration transaction",
		cfg.Migrations["CreateAnimals"].Spec.Transaction)
	expecto.True(t, "deafult is true",
		*cfg.Migrations["CreateAnimals"].Spec.Transaction)
	expecto.Eq(t, "migration up",
		cfg.Migrations["CreateAnimals"].Spec.Run.Up.Sql,
		"CREATE TABLE animals(name varchar PRIMARY KEY)")
	expecto.Eq(t, "migration down",
		cfg.Migrations["CreateAnimals"].Spec.Run.Down.Sql,
		"DROP TABLE animals")

	// migration set
	expecto.Eq(t, "migration set name",
		cfg.MigrationSets["public"].Metadata.Name,
		"public")
	expecto.Eq(t, "migration set migrations",
		cfg.MigrationSets["public"].Spec.Migrations,
		[]string{"CreateAnimals"})
}

func TestParser_WithDuplicatedRefNames(t *testing.T) {
	_, dir, cleanup := expecto.TempFS(
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

	fs := config.NewFilesystem(dir)
	parser := config.NewParser(fs)
	parsed, err := parser.LoadAndParse()
	expecto.NoErr(t, "parsing config", err)

	builder := NewBuilder(fs, parsed, plugins.New())
	_, errs := builder.BuildConfig()
	expecto.NotNil(t, "build errors", errs)
	expecto.Eq(t, "number of errors",
		len(errs.Errors), 1)
	expecto.ErrContains(t, "duplicate migration", errs.Errors[0], "adding Migration: migration `CreateAnimals` already exists")
}

func TestParser_MigrationSetWithDuplicatedRefName(t *testing.T) {
	_, dir, cleanup := expecto.TempFS(
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

	fs := config.NewFilesystem(dir)
	parser := config.NewParser(fs)
	parsed, err := parser.LoadAndParse()

	expecto.NoErr(t, "parsing config", err)

	builder := NewBuilder(fs, parsed, plugins.New())
	_, errs := builder.BuildConfig()
	expecto.NotNil(t, "build errors", errs)
	expecto.Eq(t, "number of errors",
		len(errs.Errors), 1)
	expecto.ErrContains(t, "duplicate migration set", errs.Errors[0], "adding MigrationSet: migration set `public` already exists")
}

func TestParser_MigrationSetWithMissingMigrationReference(t *testing.T) {
	_, dir, cleanup := expecto.TempFS(
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

	fs := config.NewFilesystem(dir)
	parser := config.NewParser(fs)
	parsed, err := parser.LoadAndParse()
	expecto.NoErr(t, "parsing config", err)

	builder := NewBuilder(fs, parsed, plugins.New())
	_, errs := builder.BuildConfig()
	expecto.NotNil(t, "build errors", errs)
	expecto.Eq(t, "number of errors",
		len(errs.Errors), 1)
	expecto.ErrContains(t, "missing migration", errs.Errors[0],
		"invalid reference: Migration `DoesNotExist` (referenced by MigrationSet `public`) does not exist")
}

func TestParser_MigrationsDefinedInSQLFiles(t *testing.T) {
	_, dir, cleanup := expecto.TempFS(
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
      file: "animals/up.sql"
    down:
      file: "animals/down.sql"
`,
		"src/animals/up.sql",
		`CREATE TABLE animals(name varchar PRIMARY KEY)`,
		"src/animals/down.sql",
		`DROP TABLE animals`,
	)
	defer cleanup()

	fs := config.NewFilesystem(dir)
	parser := config.NewParser(fs)
	parsed, err := parser.LoadAndParse()

	expecto.NoErr(t, "parsing config", err)
	builder := NewBuilder(fs, parsed, plugins.New())
	_, errs := builder.BuildConfig()
	expecto.NotNil(t, "build errors", errs)
	expecto.Eq(t, "number of errors", len(errs.Errors), 0)
}

func TestParser_MigrationsDefinedInSQLFilesThatAreMissing(t *testing.T) {
	_, dir, cleanup := expecto.TempFS(
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
      file: "animals/up.sql"
    down:
      file: "animals/down.sql"
`)
	defer cleanup()

	fs := config.NewFilesystem(dir)
	parser := config.NewParser(fs)
	parsed, err := parser.LoadAndParse()

	expecto.NoErr(t, "parsing config", err)
	builder := NewBuilder(fs, parsed, plugins.New())
	_, errs := builder.BuildConfig()
	expecto.NotNil(t, "build errors", errs)
	expecto.Eq(t, "number of errors",
		len(errs.Errors), 1)
	expecto.ErrContains(t, "missing migration", errs.Errors[0],
		"io error: Migration(up) `CreateAnimals` cannot load file `"+dir+"/src/animals/up.sql`")
}

func TestParser_MigrationSetValidation(t *testing.T) {
	_, dir, cleanup := expecto.TempFS(
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
spec:
  namespace:
    schema: public
    prefix: "public_"
  migrations:
    - CreateAnimals
`)
	defer cleanup()

	fs := config.NewFilesystem(dir)
	parser := config.NewParser(fs)
	parsed, err := parser.LoadAndParse()
	expecto.NoErr(t, "parsing config", err)

	builder := NewBuilder(fs, parsed, plugins.New())
	_, errs := builder.BuildConfig()
	expecto.NotNil(t, "build errors", errs)
	expecto.Eq(t, "number of errors", len(errs.Errors), 1)

	expecto.ErrContains(t, "migration set validation", errs.Errors[0], "invalid spec: MigrationSet must have a name")
}

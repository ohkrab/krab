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

	should := expecto.New(t)

	fs := config.NewFilesystem(dir)
	parser := config.NewParser(fs)
	parsed, err := parser.LoadAndParse()
	should.NoErr("parsing config", err)

	builder := NewBuilder(fs, parsed, plugins.New())
	cfg, errs := builder.BuildConfig()
	should.NotNil("build errors", errs)
	should.Eq("number of errors", len(errs.Errors), 0)

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

	should := expecto.New(t)
	fs := config.NewFilesystem(dir)
	parser := config.NewParser(fs)
	parsed, err := parser.LoadAndParse()

	should.NoErr("parsing config", err)

	builder := NewBuilder(fs, parsed, plugins.New())
	_, errs := builder.BuildConfig()
	should.NotNil("build errors", errs)
	should.Eq("number of errors",
		len(errs.Errors), 1)
	should.ErrContains("duplicate migration", errs.Errors[0], "adding Migration: migration `CreateAnimals` already exists")
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

	should := expecto.New(t)
	fs := config.NewFilesystem(dir)
	parser := config.NewParser(fs)
	parsed, err := parser.LoadAndParse()

	should.NoErr("parsing config", err)

	builder := NewBuilder(fs, parsed, plugins.New())
	_, errs := builder.BuildConfig()
	should.NotNil("build errors", errs)
	should.Eq("number of errors",
		len(errs.Errors), 1)
	should.ErrContains("duplicate migration set", errs.Errors[0], "adding MigrationSet: migration set `public` already exists")
}

func TestParser_MigrationSetWithMissingMigrationReference(t *testing.T) {
	should := expecto.New(t)

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
	should.NoErr("parsing config", err)

	builder := NewBuilder(fs, parsed, plugins.New())
	_, errs := builder.BuildConfig()
	should.NotNil("build errors", errs)
	should.Eq("number of errors",
		len(errs.Errors), 1)
	should.ErrContains("missing migration", errs.Errors[0],
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

	should := expecto.New(t)
	fs := config.NewFilesystem(dir)
	parser := config.NewParser(fs)
	parsed, err := parser.LoadAndParse()

	should.NoErr("parsing config", err)
	builder := NewBuilder(fs, parsed, plugins.New())
	_, errs := builder.BuildConfig()
	should.NotNil("build errors", errs)
	should.Eq("number of errors", len(errs.Errors), 0)
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

	should := expecto.New(t)
	fs := config.NewFilesystem(dir)
	parser := config.NewParser(fs)
	parsed, err := parser.LoadAndParse()

	should.NoErr("parsing config", err)
	builder := NewBuilder(fs, parsed, plugins.New())
	_, errs := builder.BuildConfig()
	should.NotNil("build errors", errs)
	should.Eq("number of errors",
		len(errs.Errors), 1)
	should.ErrContains("missing migration", errs.Errors[0],
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

	should := expecto.New(t)

	fs := config.NewFilesystem(dir)
	parser := config.NewParser(fs)
	parsed, err := parser.LoadAndParse()
	should.NoErr("parsing config", err)

	builder := NewBuilder(fs, parsed, plugins.New())
	_, errs := builder.BuildConfig()
	should.NotNil("build errors", errs)
	should.Eq("number of errors", len(errs.Errors), 1)

	should.ErrContains("migration set validation", errs.Errors[0], "invalid spec: MigrationSet must have a name")
}

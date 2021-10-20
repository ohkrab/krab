package krab

import (
	"context"
	"testing"

	_ "github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
	"github.com/ohkrab/krab/cli"
	"github.com/stretchr/testify/assert"
)

func TestActionMigrateDown(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	withPg(t, func(db *sqlx.DB) {
		// setup
		set := &MigrationSet{
			RefName: "public",
			Migrations: []*Migration{
				{
					Version: "v1",
					Up: MigrationUp{
						SQL: `CREATE TABLE animals(name VARCHAR)`,
					},
					Down: MigrationDown{
						SQL: `DROP TABLE animals`,
					},
				},
				{
					Version: "v2",
					Up: MigrationUp{
						SQL: `ALTER TABLE animals ADD COLUMN emoji VARCHAR`,
					},
					Down: MigrationDown{
						SQL: `ALTER TABLE animals DROP COLUMN emoji`,
					},
				},
			},
		}
		set.InitDefaults()

		err := (&ActionMigrateUp{Set: set}).Do(ctx, db, emptyTemplates(), cli.NullUI())
		assert.NoError(err, "Up migration should pass")

		_, err = db.ExecContext(ctx, "INSERT INTO animals(name, emoji) VALUES('Elephant', '🐘')")
		assert.NoError(err, "Elephant must be inserted")

		// state before
		schema, _ := NewSchemaMigrationTable("public").SelectAll(ctx, db)
		if assert.Equal(2, len(schema)) {
			assert.Equal("v1", schema[0].Version)
			assert.Equal("v2", schema[1].Version)
		}

		rowsBefore, err := db.QueryxContext(ctx, "SELECT * FROM animals")
		assert.NoError(err, "Animals must be fetched")
		defer rowsBefore.Close()

		colsBefore, _ := rowsBefore.Columns()
		assert.ElementsMatch([]string{"name", "emoji"}, colsBefore, "Columns must match")

		animals := sqlxRowsMapScan(rowsBefore)
		if assert.Equal(1, len(animals)) {
			assert.Equal("Elephant", animals[0]["name"])
			assert.Equal("🐘", animals[0]["emoji"])
		}

		// action
		action := &ActionMigrateDown{Set: set, DownMigration: SchemaMigration{"v2"}}
		err = action.Do(ctx, db, emptyTemplates())
		assert.NoError(err, "Action must succeed")

		// state after
		schema, _ = NewSchemaMigrationTable("public").SelectAll(ctx, db)
		assert.Equal(1, len(schema))
		assert.Equal("v1", schema[0].Version)

		rowsAfter, err := db.QueryxContext(ctx, "SELECT * FROM animals")
		assert.NoError(err, "Animals after emoji revert must be fetched")
		defer rowsAfter.Close()

		colsAfter, _ := rowsAfter.Columns()
		assert.ElementsMatch([]string{"name"}, colsAfter, "Only single column should exist")

		animals = sqlxRowsMapScan(rowsAfter)
		assert.Equal(1, len(animals))
		assert.Equal("Elephant", animals[0]["name"])
		assert.Nil(animals[0]["emoji"])
	})
}

func TestActionMigrateDownOnError(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	withPg(t, func(db *sqlx.DB) {
		// setup
		set := &MigrationSet{
			Migrations: []*Migration{
				{
					Version: "v1",
					Up: MigrationUp{
						SQL: `CREATE TABLE animals(name VARCHAR)`,
					},
					Down: MigrationDown{
						SQL: `DROP TABLE animals`,
					},
				},
				{
					Version: "v2",
					Up: MigrationUp{
						SQL: `ALTER TABLE animals ADD COLUMN emoji VARCHAR`,
					},
					Down: MigrationDown{
						SQL: `ALTER TABLE animals DROP COLUMN emoji; ALTER TABLE animals DROP COLUMN habitat`,
					},
				},
			},
		}
		set.InitDefaults()

		err := (&ActionMigrateUp{Set: set}).Do(ctx, db, emptyTemplates(), cli.NullUI())
		assert.NoError(err, "Up migration should pass")

		_, err = db.ExecContext(ctx, "INSERT INTO animals(name, emoji) VALUES('Elephant', '🐘')")
		assert.NoError(err, "Elephant must be inserted")

		// state before
		schema, _ := NewSchemaMigrationTable("public").SelectAll(ctx, db)
		assert.Equal(len(schema), 2)
		assert.Equal("v1", schema[0].Version)
		assert.Equal("v2", schema[1].Version)

		// action
		action := &ActionMigrateDown{Set: set, DownMigration: SchemaMigration{"v2"}}
		err = action.Do(ctx, db, emptyTemplates())
		assert.Error(err, "Migration should fail")
		assert.Contains(
			err.Error(),
			`column "habitat" of relation "animals" does not exist`,
		)

		// state after
		schema, _ = NewSchemaMigrationTable("public").SelectAll(ctx, db)
		assert.Equal(2, len(schema))
		assert.Equal("v1", schema[0].Version)
		assert.Equal("v2", schema[1].Version, "Schema information should remain untouched")

		rowsAfter, err := db.QueryxContext(ctx, "SELECT * FROM animals")
		assert.NoError(err, "Animals must be fetched")
		defer rowsAfter.Close()

		colsBefore, _ := rowsAfter.Columns()
		assert.ElementsMatch([]string{"name", "emoji"}, colsBefore, "Columns must match")
	})
}

func TestActionMigrateDownWhenSchemaDoesNotExist(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	withPg(t, func(db *sqlx.DB) {
		// setup
		set := &MigrationSet{
			Migrations: []*Migration{
				{
					Version: "v1",
					Up: MigrationUp{
						SQL: `CREATE TABLE animals(name VARCHAR)`,
					},
					Down: MigrationDown{
						SQL: `DROP TABLE animals`,
					},
				},
				{
					Version: "v2",
					Up: MigrationUp{
						SQL: `ALTER TABLE animals ADD COLUMN emoji VARCHAR`,
					},
					Down: MigrationDown{
						SQL: `ALTER TABLE animals DROP COLUMN emoji`,
					},
				},
			},
		}
		set.InitDefaults()

		err := (&ActionMigrateUp{Set: set}).Do(ctx, db, emptyTemplates(), cli.NullUI())
		assert.NoError(err, "Up migration should pass")

		// state before action 1
		schema, _ := NewSchemaMigrationTable("public").SelectAll(ctx, db)
		assert.Equal(2, len(schema))
		assert.Equal("v1", schema[0].Version)
		assert.Equal("v2", schema[1].Version)

		// action 1
		action_1 := &ActionMigrateDown{Set: set, DownMigration: SchemaMigration{"v2"}}
		err = action_1.Do(ctx, db, emptyTemplates())
		assert.NoError(err, "Migrate down should pass")

		// state after action 1
		schema, _ = NewSchemaMigrationTable("public").SelectAll(ctx, db)
		assert.Equal(1, len(schema))
		assert.Equal("v1", schema[0].Version)

		// prepare data for action 2
		_, err = db.ExecContext(ctx, "INSERT INTO animals VALUES('Crab')")
		assert.NoError(err, "Crab must be inserted")
		rowsAfter, err := db.QueryxContext(ctx, "SELECT * FROM animals")
		assert.NoError(err, "Animals exist")
		defer rowsAfter.Close()

		animals := sqlxRowsMapScan(rowsAfter)
		assert.Equal(1, len(schema))
		assert.Equal("Crab", animals[0]["name"])

		// action 2
		action_2 := &ActionMigrateDown{Set: set, DownMigration: SchemaMigration{"v2"}}
		err = action_2.Do(ctx, db, emptyTemplates())
		assert.Error(err, "Second migrate down should fail")
		assert.Contains(
			err.Error(),
			`Migration has not been run yet, nothing to rollback`,
		)
	})
}

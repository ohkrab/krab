package krab

import (
	"context"
	"strings"
	"testing"

	"github.com/franela/goblin"
	_ "github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
	"github.com/ohkrab/krab/cli"
)

func Test_ActionMigrateDown(t *testing.T) {
	withPg(t, func(db *sqlx.DB) {
		g := goblin.Goblin(t)
		ctx := context.Background()

		g.Describe("Running migrate down action", func() {
			g.AfterEach(func() {
				cleanDb(db)
			})

			g.It("Migration passess successfully", func() {
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

				err := (&ActionMigrateUp{Set: set}).Do(ctx, db, cli.NullUI())
				g.Assert(err).IsNil("Up migration should pass")

				_, err = db.ExecContext(ctx, "INSERT INTO animals(name, emoji) VALUES('Elephant', '🐘')")
				g.Assert(err).IsNil("Elephant must be inserted")

				// state before
				schema, _ := SchemaMigrationSelectAll(ctx, db)
				g.Assert(len(schema)).Eql(2)
				g.Assert(schema[0].Version).Eql("v1")
				g.Assert(schema[1].Version).Eql("v2")

				rowsBefore, err := db.QueryxContext(ctx, "SELECT * FROM animals")
				g.Assert(err).IsNil("Animals must be fetched")
				defer rowsBefore.Close()

				colsBefore, _ := rowsBefore.Columns()
				g.Assert(colsBefore).Eql([]string{"name", "emoji"}, "Columns must match")

				animals := sqlxRowsMapScan(rowsBefore)
				g.Assert(len(animals)).Eql(1)
				g.Assert(animals[0]["name"]).Eql("Elephant")
				g.Assert(animals[0]["emoji"]).Eql("🐘")

				// action
				action := &ActionMigrateDown{Set: set, DownMigration: SchemaMigration{"v2"}}
				err = action.Do(ctx, db)
				g.Assert(err).IsNil("Action must succeed", err)

				// state after
				schema, _ = SchemaMigrationSelectAll(ctx, db)
				g.Assert(len(schema)).Eql(1)
				g.Assert(schema[0].Version).Eql("v1")

				rowsAfter, err := db.QueryxContext(ctx, "SELECT * FROM animals")
				g.Assert(err).IsNil("Animals after emoji revert must be fetched")
				defer rowsAfter.Close()

				colsAfter, _ := rowsAfter.Columns()
				g.Assert(colsAfter).Eql([]string{"name"}, "Only single column should exist")

				animals = sqlxRowsMapScan(rowsAfter)
				g.Assert(len(animals)).Eql(1)
				g.Assert(animals[0]["name"]).Eql("Elephant")
				g.Assert(animals[0]["emoji"]).Eql(nil)
			})

			g.It("Migration is not saved when error occurred", func() {
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

				err := (&ActionMigrateUp{Set: set}).Do(ctx, db, cli.NullUI())
				g.Assert(err).IsNil("Up migration should pass")

				_, err = db.ExecContext(ctx, "INSERT INTO animals(name, emoji) VALUES('Elephant', '🐘')")
				g.Assert(err).IsNil("Elephant must be inserted")

				// state before
				schema, _ := SchemaMigrationSelectAll(ctx, db)
				g.Assert(len(schema)).Eql(2)
				g.Assert(schema[0].Version).Eql("v1")
				g.Assert(schema[1].Version).Eql("v2")

				// action
				action := &ActionMigrateDown{Set: set, DownMigration: SchemaMigration{"v2"}}
				err = action.Do(ctx, db)
				g.Assert(err).IsNotNil("Migration should fail")
				g.Assert(
					strings.Contains(
						err.Error(),
						`column "habitat" of relation "animals" does not exist`,
					),
				).Eql(true, err)

				// state after
				schema, err = SchemaMigrationSelectAll(ctx, db)
				g.Assert(len(schema)).Eql(2)
				g.Assert(schema[0].Version).Eql("v1")
				g.Assert(schema[1].Version).Eql("v2", "Schema information should remain untouched")

				rowsAfter, err := db.QueryxContext(ctx, "SELECT * FROM animals")
				g.Assert(err).IsNil("Animals must be fetched")
				defer rowsAfter.Close()

				colsBefore, _ := rowsAfter.Columns()
				g.Assert(colsBefore).Eql([]string{"name", "emoji"}, "Columns must match")
			})

			g.It("Migration is executed when does not exist in schema", func() {
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

				err := (&ActionMigrateUp{Set: set}).Do(ctx, db, cli.NullUI())
				g.Assert(err).IsNil("Up migration should pass")

				// state before action 1
				schema, _ := SchemaMigrationSelectAll(ctx, db)
				g.Assert(len(schema)).Eql(2)
				g.Assert(schema[0].Version).Eql("v1")
				g.Assert(schema[1].Version).Eql("v2")

				// action 1
				action_1 := &ActionMigrateDown{Set: set, DownMigration: SchemaMigration{"v2"}}
				err = action_1.Do(ctx, db)
				g.Assert(err).IsNil("Migrate down should pass")

				// state after action 1
				schema, err = SchemaMigrationSelectAll(ctx, db)
				g.Assert(len(schema)).Eql(1)
				g.Assert(schema[0].Version).Eql("v1")

				// prepare data for action 2
				_, err = db.ExecContext(ctx, "INSERT INTO animals VALUES('Crab')")
				g.Assert(err).IsNil("Crab must be inserted")
				rowsAfter, err := db.QueryxContext(ctx, "SELECT * FROM animals")
				g.Assert(err).IsNil("Animals exist")
				defer rowsAfter.Close()

				animals := sqlxRowsMapScan(rowsAfter)
				g.Assert(len(animals)).Eql(1)
				g.Assert(animals[0]["name"]).Eql("Crab")

				// action 2
				action_2 := &ActionMigrateDown{Set: set, DownMigration: SchemaMigration{"v2"}}
				err = action_2.Do(ctx, db)
				g.Assert(err).IsNotNil("Second migrate down should fail")
				g.Assert(
					strings.Contains(
						err.Error(),
						`Migration has not been run yet, nothing to rollback`,
					),
				).Eql(true, err)
			})
		})
	})
}

package krab

import (
	"context"
	"testing"

	"github.com/franela/goblin"
	_ "github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
	"github.com/ohkrab/krab/cli"
)

func Test_ActionMigrateDownTransactions(t *testing.T) {
	withPg(t, func(db *sqlx.DB) {
		g := goblin.Goblin(t)
		ctx := context.Background()

		g.Describe("Running migrate down action with concurrent operation", func() {
			g.AfterEach(func() {
				cleanDb(db)
			})

			g.It("Migration passess successfully", func() {
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
					},
				}
				set.InitDefaults()

				err := (&ActionMigrateUp{Set: set}).Do(ctx, db, cli.NullUI())
				g.Assert(err).IsNil("First migration should pass")

				inTransaction := false
				set.Migrations = []*Migration{
					{
						Transaction: &inTransaction,
						Version:     "v2",
						Up: MigrationUp{
							SQL: `CREATE INDEX idx ON animals(name)`,
						},
						Down: MigrationDown{
							SQL: `DROP INDEX CONCURRENTLY idx`,
						},
					},
				}

				err = (&ActionMigrateUp{Set: set}).Do(ctx, db, cli.NullUI())
				g.Assert(err).IsNil("Second migration should pass")

				schema, err := SchemaMigrationTable{}.SelectAll(ctx, db)
				if err != nil {
					t.Error("Fetching migrations failed", err)
					return
				}

				g.Assert(len(schema)).Eql(2)
				g.Assert(schema[0].Version).Eql("v1")
				g.Assert(schema[1].Version).Eql("v2")

				err = (&ActionMigrateDown{Set: set, DownMigration: SchemaMigration{"v2"}}).Do(ctx, db)
				g.Assert(err).IsNil("Rollback of v2 migration should pass")

				schema, err = SchemaMigrationTable{}.SelectAll(ctx, db)
				if err != nil {
					t.Error("Fetching migrations failed", err)
					return
				}

				g.Assert(len(schema)).Eql(1)
				g.Assert(schema[0].Version).Eql("v1")
			})
		})
	})
}

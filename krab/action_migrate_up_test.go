package krab

import (
	"context"
	"strings"
	"testing"

	"github.com/franela/goblin"
	_ "github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
)

func Test_ActionMigrateUp(t *testing.T) {
	withPg(t, func(db *sqlx.DB) {
		g := goblin.Goblin(t)
		ctx := context.Background()

		g.Describe("Running migrate up action", func() {
			g.AfterEach(func() {
				cleanDb(db)
			})

			g.It("Migration passess successfully", func() {
				action := &ActionMigrateUp{
					Set: &MigrationSet{
						Migrations: []*Migration{
							{
								Version: "v1",
								Up: MigrationUp{
									SQL: `SELECT 1`,
								},
							},
						},
					},
				}

				err := action.Do(ctx, db)
				if err != nil {
					t.Error("Migration error:", err)
					return
				}

				schema, err := SchemaMigrationSelectAll(ctx, db)
				if err != nil {
					t.Error("Fetching migrations failed", err)
					return
				}

				g.Assert(len(schema)).Eql(1)
				g.Assert(schema[0].Version).Eql("v1")
			})

			g.It("Migration is not saved when error occurred", func() {
				SchemaMigrationInit(ctx, db)

				action := &ActionMigrateUp{
					Set: &MigrationSet{
						Migrations: []*Migration{
							{
								Version: "v1",
								Up: MigrationUp{
									SQL: `SELECT invalid`,
								},
							},
						},
					},
				}

				err := action.Do(ctx, db)

				g.Assert(err).IsNotNil("Invalid migration should return error")
				g.Assert(strings.Contains(
					err.Error(),
					`column "invalid" does not exist`,
				)).Eql(true)

				schema, err := SchemaMigrationSelectAll(ctx, db)
				if err != nil {
					t.Error("Fetching migrations failed", err)
					return
				}

				g.Assert(len(schema)).Eql(0)
			})
		})
	})
}

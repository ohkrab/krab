package krab

import (
	"context"
	"testing"

	"github.com/franela/goblin"
	"github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4"
)

func Test_ActionMigrateUp(t *testing.T) {
	g := goblin.Goblin(t)
	ctx := context.Background()

	withPg(t, func(conn *pgx.Conn) {
		g.Describe("Running migrate up action", func() {
			g.It("Migration passess successfuly", func() {
				SchemaMigrationTruncate(ctx, conn)

				action := &ActionMigrateUp{
					Set: &MigrationSet{
						Migrations: []*Migration{
							{
								RefName: "v1",
								Up: MigrationUp{
									Sql: `SELECT 1`,
								},
							},
						},
					},
				}

				err := action.Run(ctx, conn)
				if err != nil {
					t.Error("Test migration error:", err)
					return
				}

				schema, err := SchemaMigrationSelectAll(ctx, conn)
				if err != nil {
					t.Error("Fetching migrations failed", err)
					return
				}

				g.Assert(len(schema)).Eql(1)
				g.Assert(schema[0].Version).Eql("v1")
			})

			g.Xit("Migration is not saved when error occured", func() {
				SchemaMigrationInit(ctx, conn)

				action := &ActionMigrateUp{
					Set: &MigrationSet{
						Migrations: []*Migration{
							{
								RefName: "v1",
								Up: MigrationUp{
									Sql: `SELECT invalid`,
								},
							},
						},
					},
				}

				err := action.Run(ctx, conn)
				g.Assert(err.Error()).Eql(`pq: column "invalid" does not exist`)

				schema, err := SchemaMigrationSelectAll(ctx, conn)
				if err != nil {
					t.Error("Fetching migrations failed", err)
					return
				}

				g.Assert(len(schema)).Eql(0)
			})
		})
	})
}

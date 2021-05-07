package krab

import (
	"context"
	"testing"

	epg "github.com/fergusstrange/embedded-postgres"
	"github.com/franela/goblin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Test_ActionMigrateUp(t *testing.T) {
	g := goblin.Goblin(t)
	ctx := context.Background()

	withPg(t, func(db *sqlx.DB) {
		g.Describe("#Run", func() {
			g.It("Migration passess successfuly", func() {
				action := &ActionMigrateUp{
					db: db,
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

				err := action.Run(ctx)
				if err != nil {
					t.Error("Migration error", err)
					return
				}

				var schema []SchemaInfo
				db.Select(&schema, "SELECT * FROM schema_info")

				g.Assert(len(schema)).Eql(1)
				g.Assert(schema[0].Version).Eql("v1")
			})
		})
	})
}

func withPg(t *testing.T, f func(db *sqlx.DB)) {
	database := epg.NewDatabase()

	if err := database.Start(); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := database.Stop(); err != nil {
			t.Fatal(err)
		}
	}()

	db, err := sqlx.Connect(
		"postgres",
		"host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable",
	)

	if err != nil {
		t.Fatal(err)
	}
	f(db)
}

package krab

import (
	"context"
	"testing"

	_ "github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
	"github.com/ohkrab/krab/cli"
	"github.com/stretchr/testify/assert"
)

func TestActionMigrateUpHooks(t *testing.T) {
	assert := assert.New(t)

	withPg(t, func(db *sqlx.DB) {
		defer cleanDb(db)
		ctx := context.Background()

		db.MustExec("CREATE SCHEMA tenants")
		defer db.MustExec("DROP SCHEMA tenants CASCADE")

		set := &MigrationSet{
			Hooks: &Hooks{
				Before: "SET search_path TO tenants",
			},
			RefName: "tenants",
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
		assert.NoError(err, "First migration should pass")

		schema, err := SchemaMigrationTable{"public.schema_migrations"}.SelectAll(ctx, db)
		assert.Equal(0, len(schema))
		if assert.Error(err) {
			assert.Contains(err.Error(), `relation "public.schema_migrations" does not exist`)
		}

		schema, err = SchemaMigrationTable{"tenants.schema_migrations"}.SelectAll(ctx, db)
		assert.NoError(err, "Fetching migrations from tenant schema should be successful")

		if assert.Equal(1, len(schema)) {
			assert.Equal("v1", schema[0].Version)
		}
	})
}

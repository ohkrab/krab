package krab

import (
	"context"
	"testing"

	_ "github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
	"github.com/ohkrab/krab/cli"
	"github.com/stretchr/testify/assert"
)

func TestActionMigrateUp(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	withPg(t, func(db *sqlx.DB) {
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
		action.Set.InitDefaults()

		err := action.Do(ctx, db, cli.NullUI())
		assert.NoError(err)

		schema, err := SchemaMigrationTable{}.SelectAll(ctx, db)
		assert.NoError(err)

		assert.Equal(1, len(schema))
		assert.Equal("v1", schema[0].Version)
	})
}

func TestActionMigrateUpWithError(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	withPg(t, func(db *sqlx.DB) {
		SchemaMigrationTable{}.Init(ctx, db)

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
		action.Set.InitDefaults()

		err := action.Do(ctx, db, cli.NullUI())

		assert.Error(err)
		assert.Contains(
			err.Error(),
			`column "invalid" does not exist`,
		)

		schema, err := SchemaMigrationTable{}.SelectAll(ctx, db)
		assert.Equal(len(schema), 0)
	})
}

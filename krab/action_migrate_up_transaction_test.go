package krab

import (
	"context"
	"testing"

	_ "github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
	"github.com/ohkrab/krab/cli"
	"github.com/stretchr/testify/assert"
)

func TestActionMigrateUpTransactions(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	withPg(t, func(db *sqlx.DB) {
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
		assert.NoError(err, "First migration should pass")

		inTransaction := false
		set.Migrations = []*Migration{
			{
				Transaction: &inTransaction,
				Version:     "v2",
				Up: MigrationUp{
					SQL: `CREATE INDEX CONCURRENTLY idx ON animals(name)`,
				},
				Down: MigrationDown{
					SQL: `DROP INDEX animals`,
				},
			},
		}

		err = (&ActionMigrateUp{Set: set}).Do(ctx, db, cli.NullUI())
		assert.NoError(err, "Second migration should pass")

		schema, err := NewSchemaMigrationTable("public").SelectAll(ctx, db)
		assert.NoError(err, "Fetching migrations failed")

		assert.Equal(2, len(schema))
		assert.Equal("v1", schema[0].Version)
		assert.Equal("v2", schema[1].Version)
	})
}

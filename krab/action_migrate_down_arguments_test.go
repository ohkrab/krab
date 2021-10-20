package krab

import (
	"context"
	"testing"

	_ "github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
	"github.com/ohkrab/krab/cli"
	"github.com/ohkrab/krab/tpls"
	"github.com/stretchr/testify/assert"
)

func TestActionMigrateDownArguments(t *testing.T) {
	assert := assert.New(t)

	withPg(t, func(db *sqlx.DB) {
		ctx := context.Background()

		set := createMigrationSet("tenants",
			"v1",
			`CREATE TABLE animals(name VARCHAR)`,
			`DROP TABLE animals`,
		)
		set.Schema = "{{.Args.schema}}"
		set.Arguments = &Arguments{
			Args: []*Argument{
				{
					Name: "schema",
					Type: "string",
				},
			},
		}

		templates := tpls.New(map[string]interface{}{
			"schema": "custom",
		})
		err := (&ActionMigrateUp{Set: set}).Do(ctx, db, templates, cli.NullUI())
		assert.NoError(err, "First migration should pass")

		schema, err := SchemaMigrationTable{"public.schema_migrations"}.SelectAll(ctx, db)
		assert.Equal(0, len(schema))
		if assert.Error(err) {
			assert.Contains(err.Error(), `relation "public.schema_migrations" does not exist`)
		}

		schema, err = SchemaMigrationTable{"custom.schema_migrations"}.SelectAll(ctx, db)
		assert.NoError(err, "Fetching migrations from tenant schema should be successful")

		if assert.Equal(1, len(schema)) {
			assert.Equal("v1", schema[0].Version)
		}
		err = (&ActionMigrateDown{Set: set, DownMigration: SchemaMigration{"v1"}}).Do(ctx, db, templates)
		assert.NoError(err, "Down migration should pass")

		schema, err = SchemaMigrationTable{"custom.schema_migrations"}.SelectAll(ctx, db)
		assert.NoError(err, "Fetching migrations from tenant schema should be successful")

		assert.Equal(0, len(schema))
	})
}

package krab

import (
	"context"
	"testing"

	_ "github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
	"github.com/ohkrab/krab/cli"
	"github.com/stretchr/testify/assert"
)

func TestActionMigrateDownHooks(t *testing.T) {
	assert := assert.New(t)

	withPg(t, func(db *sqlx.DB) {
		ctx := context.Background()

		set := createMigrationSet("tenants",
			"v1",
			`CREATE TABLE animals(name VARCHAR)`,
			`DROP TABLE animals`,
			"v2",
			`ALTER TABLE animals ADD COLUMN emoji VARCHAR`,
			`ALTER TABLE animals DROP COLUMN emoji`,
		)
		set.Schema = "tenants"

		err := (&ActionMigrateUp{Set: set}).Do(ctx, db, emptyTemplates(), cli.NullUI())
		assert.NoError(err, "First migration should pass")

		schema, _ := SchemaMigrationTable{"tenants.schema_migrations"}.SelectAll(ctx, db)
		if assert.Equal(2, len(schema)) {
			assert.Equal("v1", schema[0].Version)
			assert.Equal("v2", schema[1].Version)
		}

		err = (&ActionMigrateDown{Set: set, DownMigration: SchemaMigration{"v2"}}).Do(ctx, db, emptyTemplates())
		assert.NoError(err, "Rollback migration should pass")

		schema, _ = SchemaMigrationTable{"tenants.schema_migrations"}.SelectAll(ctx, db)
		if assert.Equal(1, len(schema)) {
			assert.Equal("v1", schema[0].Version)
		}
	})
}

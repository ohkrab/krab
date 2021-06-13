package krab

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// SchemaMigration represents a single row from migrations table.
type SchemaMigration struct {
	Version string `db:"version"`
}

// SchemaMigrationInit creates a migrations table.
func SchemaMigrationInit(ctx context.Context, db *sqlx.DB) error {
	_, err := db.ExecContext(ctx, fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s(version varchar PRIMARY KEY)",
		pq.QuoteIdentifier(DefaultMigrationsTableName),
	))
	return err
}

// SchemaMigrationSelectAll fetches all migrations from a database.
func SchemaMigrationSelectAll(ctx context.Context, db *sqlx.DB) ([]SchemaMigration, error) {
	var schema []SchemaMigration
	err := db.SelectContext(
		ctx,
		&schema,
		fmt.Sprintf("SELECT version FROM %s ORDER BY 1", pq.QuoteIdentifier(DefaultMigrationsTableName)),
	)
	return schema, err
}

// SchemaMigrationInsert saves migration to a database.
func SchemaMigrationInsert(ctx context.Context, db *sqlx.DB, refName string) error {
	_, err := db.ExecContext(ctx, fmt.Sprintf(
		"INSERT INTO %s(version) VALUES ($1)",
		pq.QuoteIdentifier(DefaultMigrationsTableName),
	),
		refName,
	)
	return err
}

// SchemaMigrationFilterPending removes `refsInDb` migrations from `all` and return new slice with pending ones only.
func SchemaMigrationFilterPending(all []*Migration, refsInDb []SchemaMigration) []*Migration {
	pendingMigrations := make([]*Migration, 0)

	for _, migration := range all {
		var found *Migration
		for _, ref := range refsInDb {
			if migration.RefName == ref.Version {
				found = migration
				break
			}
		}

		if found == nil {
			pendingMigrations = append(pendingMigrations, migration)
		}
	}

	return pendingMigrations
}

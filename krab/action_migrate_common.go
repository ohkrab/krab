package krab

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/ohkrab/krab/krabdb"
)

// SchemaMigration represents a single row from migrations table.
type SchemaMigration struct {
	Version string `db:"version"`
}

// SchemaMigrationInit creates a migrations table.
func SchemaMigrationInit(ctx context.Context, db sqlx.ExecerContext) error {
	_, err := db.ExecContext(ctx, fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s(version varchar PRIMARY KEY)",
		krabdb.QuoteIdent(defaultMigrationsTableName),
	))
	return err
}

// SchemaMigrationTruncate truncates migrations table.
func SchemaMigrationTruncate(ctx context.Context, db sqlx.ExecerContext) error {
	_, err := db.ExecContext(ctx, fmt.Sprintf(
		"TRUNCATE %s",
		krabdb.QuoteIdent(defaultMigrationsTableName),
	))
	return err
}

// SchemaMigrationExists checks if migration exists in database.
func SchemaMigrationExists(ctx context.Context, db sqlx.QueryerContext, migration SchemaMigration) (bool, error) {
	var schema []SchemaMigration
	err := sqlx.SelectContext(
		ctx,
		db,
		&schema,
		fmt.Sprintf("SELECT version FROM %s WHERE version = $1", krabdb.QuoteIdent(defaultMigrationsTableName)),
		migration.Version,
	)
	return len(schema) > 0, err
}

// SchemaMigrationSelectLastN fetches last N migrations in Z-A order.
func SchemaMigrationSelectLastN(ctx context.Context, db sqlx.QueryerContext, limit int) ([]SchemaMigration, error) {
	var schema []SchemaMigration
	err := sqlx.SelectContext(
		ctx,
		db,
		&schema,
		fmt.Sprintf("SELECT version FROM %s ORDER BY 1 DESC LIMIT %d", krabdb.QuoteIdent(defaultMigrationsTableName), limit),
	)
	return schema, err
}

// SchemaMigrationSelectAll fetches all migrations from a database.
func SchemaMigrationSelectAll(ctx context.Context, db sqlx.QueryerContext) ([]SchemaMigration, error) {
	var schema []SchemaMigration
	err := sqlx.SelectContext(
		ctx,
		db,
		&schema,
		fmt.Sprintf("SELECT version FROM %s ORDER BY 1", krabdb.QuoteIdent(defaultMigrationsTableName)),
	)
	return schema, err
}

// SchemaMigrationInsert saves migration to a database.
func SchemaMigrationInsert(ctx context.Context, db sqlx.ExecerContext, version string) error {
	_, err := db.ExecContext(
		ctx,
		fmt.Sprintf("INSERT INTO %s(version) VALUES ($1) RETURNING *", krabdb.QuoteIdent(defaultMigrationsTableName)),
		version,
	)
	return err
}

// SchemaMigrationDelete removes migration from a database.
func SchemaMigrationDelete(ctx context.Context, db sqlx.ExecerContext, version string) error {
	_, err := db.ExecContext(
		ctx,
		fmt.Sprintf("DELETE FROM %s WHERE version = $1 RETURNING *", krabdb.QuoteIdent(defaultMigrationsTableName)),
		version,
	)
	return err
}

// SchemaMigrationFilterPending removes `refsInDb` migrations from `all` and return new slice with pending ones only.
func SchemaMigrationFilterPending(all []*Migration, refsInDb []SchemaMigration) []*Migration {
	pendingMigrations := make([]*Migration, 0)

	for _, migration := range all {
		var found *Migration
		for _, ref := range refsInDb {
			if migration.Version == ref.Version {
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

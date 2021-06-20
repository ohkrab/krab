package krab

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/ohkrab/krab/krabdb"
)

// SchemaMigration represents a single row from migrations table.
type SchemaMigration struct {
	Version string `db:"version"`
}

// SchemaMigrationInit creates a migrations table.
func SchemaMigrationInit(ctx context.Context, conn *pgx.Conn) error {
	_, err := conn.Exec(ctx, fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s(version varchar PRIMARY KEY)",
		krabdb.QuoteIdent(DefaultMigrationsTableName),
	))
	return err
}

// SchemaMigrationTruncate truncates migrations table.
func SchemaMigrationTruncate(ctx context.Context, conn *pgx.Conn) error {
	_, err := conn.Exec(ctx, fmt.Sprintf(
		"TRUNCATE %s",
		krabdb.QuoteIdent(DefaultMigrationsTableName),
	))
	return err
}

// SchemaMigrationSelectAll fetches all migrations from a database.
func SchemaMigrationSelectAll(ctx context.Context, conn *pgx.Conn) ([]SchemaMigration, error) {
	var schema []SchemaMigration
	rows, err := conn.Query(
		ctx,
		fmt.Sprintf("SELECT version FROM %s ORDER BY 1", krabdb.QuoteIdent(DefaultMigrationsTableName)),
	)
	defer rows.Close()

	return schema, err
}

// SchemaMigrationInsert saves migration to a database.
func SchemaMigrationInsert(ctx context.Context, conn *pgx.Conn, refName string) error {
	_, err := conn.Exec(
		ctx,
		fmt.Sprintf("INSERT INTO %s(version) VALUES ($1) RETURNING *", krabdb.QuoteIdent(DefaultMigrationsTableName)),
		// fmt.Sprintf("INSERT INTO %s(version) VALUES ('%s')", krabdb.QuoteIdent(DefaultMigrationsTableName), refName),
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

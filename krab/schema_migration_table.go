package krab

import (
	"context"
	"fmt"
	"strings"

	"github.com/ohkrab/krab/krabdb"
)

const DefaultSchemaMigrationTableName = "schema_migrations"

type SchemaMigrationTable struct {
	Name string
}

// NewSchemaMigrationTable creates SchemaMigrationTable with default table name and specified schema.
func NewSchemaMigrationTable(schema string) SchemaMigrationTable {
	return SchemaMigrationTable{
		Name: strings.Join(
			[]string{
				schema,
				DefaultSchemaMigrationTableName,
			},
			".",
		),
	}
}

// SchemaMigration represents a single row from migrations table.
type SchemaMigration struct {
	Version string `db:"version"`
}

// Init creates a migrations table.
func (s SchemaMigrationTable) Init(ctx context.Context, db krabdb.ExecerContext) error {
	parts := strings.Split(s.Name, ".")
	if len(parts) > 1 {
		schema := parts[0]
		_, err := db.ExecContext(ctx, fmt.Sprintf(
			"CREATE SCHEMA IF NOT EXISTS %s",
			krabdb.QuoteIdent(schema),
		))
		if err != nil {
			return err
		}
	}

	_, err := db.ExecContext(ctx, fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s(version varchar PRIMARY KEY, migrated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP)",
		krabdb.QuoteIdentWithDots(s.Name),
	))

	return err
}

// Truncate truncates migrations table.
func (s SchemaMigrationTable) Truncate(ctx context.Context, db krabdb.ExecerContext) error {
	_, err := db.ExecContext(ctx, fmt.Sprintf(
		"TRUNCATE %s",
		krabdb.QuoteIdentWithDots(s.Name),
	))
	return err
}

// Exists checks if migration exists in database.
func (s SchemaMigrationTable) Exists(ctx context.Context, db krabdb.QueryerContext, migration SchemaMigration) (bool, error) {
	var schema []SchemaMigration
	err := db.SelectContext(
		ctx,
		&schema,
		fmt.Sprintf("SELECT version FROM %s WHERE version = $1", krabdb.QuoteIdentWithDots(s.Name)),
		migration.Version,
	)
	return len(schema) > 0, err
}

// SelectLastN fetches last N migrations in Z-A order.
func (s SchemaMigrationTable) SelectLastN(ctx context.Context, db krabdb.QueryerContext, limit int) ([]SchemaMigration, error) {
	var schema []SchemaMigration
	err := db.SelectContext(
		ctx,
		&schema,
		fmt.Sprintf("SELECT version FROM %s ORDER BY 1 DESC LIMIT %d", krabdb.QuoteIdentWithDots(s.Name), limit),
	)
	return schema, err
}

// SelectAll fetches all migrations from a database.
func (s SchemaMigrationTable) SelectAll(ctx context.Context, db krabdb.QueryerContext) ([]SchemaMigration, error) {
	var schema []SchemaMigration
	err := db.SelectContext(
		ctx,
		&schema,
		fmt.Sprintf("SELECT version FROM %s ORDER BY 1", krabdb.QuoteIdentWithDots(s.Name)),
	)
	return schema, err
}

// Insert saves migration to a database.
func (s SchemaMigrationTable) Insert(ctx context.Context, db krabdb.ExecerContext, version string) error {
	_, err := db.ExecContext(
		ctx,
		fmt.Sprintf("INSERT INTO %s(version) VALUES ($1) RETURNING *", krabdb.QuoteIdentWithDots(s.Name)),
		version,
	)
	return err
}

// Delete removes migration from a database.
func (s SchemaMigrationTable) Delete(ctx context.Context, db krabdb.ExecerContext, version string) error {
	_, err := db.ExecContext(
		ctx,
		fmt.Sprintf("DELETE FROM %s WHERE version = $1 RETURNING *", krabdb.QuoteIdentWithDots(s.Name)),
		version,
	)
	return err
}

// FilterPending removes `refsInDb` migrations from `all` and return new slice with pending ones only.
func (s SchemaMigrationTable) FilterPending(all []*Migration, refsInDb []SchemaMigration) []*Migration {
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

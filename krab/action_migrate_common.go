package krab

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/ohkrab/krab/krabdb"
)

type SchemaMigrationTable struct {
	Name string
}

// SchemaMigration represents a single row from migrations table.
type SchemaMigration struct {
	Version string `db:"version"`
}

// Init creates a migrations table.
func (s SchemaMigrationTable) Init(ctx context.Context, db sqlx.ExecerContext) error {
	_, err := db.ExecContext(ctx, fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s(version varchar PRIMARY KEY)",
		krabdb.QuoteIdent(s.TableName()),
	))
	return err
}

// Truncate truncates migrations table.
func (s SchemaMigrationTable) Truncate(ctx context.Context, db sqlx.ExecerContext) error {
	_, err := db.ExecContext(ctx, fmt.Sprintf(
		"TRUNCATE %s",
		krabdb.QuoteIdent(s.TableName()),
	))
	return err
}

// Exists checks if migration exists in database.
func (s SchemaMigrationTable) Exists(ctx context.Context, db sqlx.QueryerContext, migration SchemaMigration) (bool, error) {
	var schema []SchemaMigration
	err := sqlx.SelectContext(
		ctx,
		db,
		&schema,
		fmt.Sprintf("SELECT version FROM %s WHERE version = $1", krabdb.QuoteIdent(s.TableName())),
		migration.Version,
	)
	return len(schema) > 0, err
}

// SelectLastN fetches last N migrations in Z-A order.
func (s SchemaMigrationTable) SelectLastN(ctx context.Context, db sqlx.QueryerContext, limit int) ([]SchemaMigration, error) {
	var schema []SchemaMigration
	err := sqlx.SelectContext(
		ctx,
		db,
		&schema,
		fmt.Sprintf("SELECT version FROM %s ORDER BY 1 DESC LIMIT %d", krabdb.QuoteIdent(s.TableName()), limit),
	)
	return schema, err
}

// SelectAll fetches all migrations from a database.
func (s SchemaMigrationTable) SelectAll(ctx context.Context, db sqlx.QueryerContext) ([]SchemaMigration, error) {
	var schema []SchemaMigration
	err := sqlx.SelectContext(
		ctx,
		db,
		&schema,
		fmt.Sprintf("SELECT version FROM %s ORDER BY 1", krabdb.QuoteIdent(s.TableName())),
	)
	return schema, err
}

// Insert saves migration to a database.
func (s SchemaMigrationTable) Insert(ctx context.Context, db sqlx.ExecerContext, version string) error {
	_, err := db.ExecContext(
		ctx,
		fmt.Sprintf("INSERT INTO %s(version) VALUES ($1) RETURNING *", krabdb.QuoteIdent(s.TableName())),
		version,
	)
	return err
}

// Delete removes migration from a database.
func (s SchemaMigrationTable) Delete(ctx context.Context, db sqlx.ExecerContext, version string) error {
	_, err := db.ExecContext(
		ctx,
		fmt.Sprintf("DELETE FROM %s WHERE version = $1 RETURNING *", krabdb.QuoteIdent(s.TableName())),
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

func (s SchemaMigrationTable) TableName() string {
	if s.Name == "" {
		return "schema_migrations"
	}

	return s.Name
}

package krab

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/ohkrab/krab/krabdb"
	"github.com/pkg/errors"
)

type SchemaInfo struct {
	Version string `db:"version"`
}

type ActionMigrateUp struct {
	Set *MigrationSet
	db  *sqlx.DB
}

func (a *ActionMigrateUp) fetchMigrationsFromDb(ctx context.Context) ([]SchemaInfo, error) {
	var schema []SchemaInfo
	err := a.db.SelectContext(ctx, &schema, fmt.Sprintf("SELECT * FROM %s", pq.QuoteIdentifier(DefaultMigrationsTableName)))
	return schema, err
}

func (a *ActionMigrateUp) insertToSchemaInformation(ctx context.Context, refName string) error {
	_, err := a.db.ExecContext(ctx, fmt.Sprintf(
		"INSERT INTO %s(version) VALUES ($1)",
		pq.QuoteIdentifier(DefaultMigrationsTableName),
	),
		refName,
	)
	return err
}

func (a *ActionMigrateUp) createSchema(ctx context.Context) error {
	_, err := a.db.ExecContext(ctx, fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s(version varchar PRIMARY KEY)",
		pq.QuoteIdentifier(DefaultMigrationsTableName),
	))
	return err
}

func (a *ActionMigrateUp) migrate(ctx context.Context, migration *Migration) error {
	// BEGIN
	_, err := a.db.ExecContext(ctx, migration.Up.Sql)
	if err != nil {
		return err
	}

	err = a.insertToSchemaInformation(ctx, migration.RefName)
	if err != nil {
		// ROLLBACK
		return err
	}

	// COMMIT
	return nil
}

func (a *ActionMigrateUp) Run(ctx context.Context) error {
	tx, err := a.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Commit()

	ok, err := krabdb.TryAdvisoryXactLock(ctx, tx, 1)
	if err != nil {
		return err
	}

	if ok {
		err := a.createSchema(ctx)
		// TODO: ensure structure compatiblity with inner structs
		if err != nil {
			return errors.Wrap(err, "Failed to create default table for migrations")
		}

		migrationRefsInDb, err := a.fetchMigrationsFromDb(ctx)
		if err != nil {
			return err
		}

		pendingMigrations := a.findPendingMigrations(migrationRefsInDb)

		{
			tx, err := a.db.BeginTxx(ctx, nil)
			if err != nil {
				return err
			}

			for _, pending := range pendingMigrations {
				err := a.migrate(ctx, pending)
				if err != nil {
					tx.Rollback() // ignore rollback error
					return err
				}
			}

			err = tx.Commit()
			if err != nil {
				return err
			}
		}
	} else {
		return errors.New("Another migration in progress")
	}

	return nil
}

func (a *ActionMigrateUp) findPendingMigrations(refsInDb []SchemaInfo) []*Migration {
	pendingMigrations := make([]*Migration, 0)

	for _, migration := range a.Set.Migrations {
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

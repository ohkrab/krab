package krab

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
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

func (a *ActionMigrateUp) fetchMigrationsFromDb(ctx context.Context) ([]string, error) {
	return []string{"test"}, nil
}

func (a *ActionMigrateUp) execInDb(ctx context.Context, sql string) error {
	fmt.Println("Migrating", sql)
	return nil
}

func (a *ActionMigrateUp) insertToSchemaInformation(ctx context.Context, refName string) error {
	_, err := a.db.NamedExec("INSERT INTO schema_info(version) VALUES (:version)", map[string]interface{}{"version": refName})
	return err
}

func (a *ActionMigrateUp) createSchemeInformation(ctx context.Context) error {
	_, err := a.db.Exec("CREATE TABLE IF NOT EXIST schema_info(version varchar PRIMARY KEY)")
	return err
}

func (a *ActionMigrateUp) migrate(ctx context.Context, migration *Migration) error {
	// BEGIN
	err := a.execInDb(ctx, migration.Up.Sql)
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
	lock := krabdb.AdvisoryLock{Errs: make(chan error)}

	fmt.Println("{{{{")
	err := a.createSchemeInformation(ctx)
	// TODO: ensure schema_info structure compatiblity with inner structs
	if err != nil {
		return errors.Wrap(err, "Failed to create `schema_info` table for migrations")
	}
	fmt.Println("}}}}")

	lock.Lock(ctx)
	defer lock.Unlock(ctx)

	migrationRefsInDb, err := a.fetchMigrationsFromDb(ctx)
	if err != nil {
		return err
	}

	pendingMigrations := a.findPendingMigrations(migrationRefsInDb)

	for _, pending := range pendingMigrations {
		err := a.migrate(ctx, pending)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *ActionMigrateUp) findPendingMigrations(refsInDb []string) []*Migration {
	pendingMigrations := make([]*Migration, 0)

	for _, migration := range a.Set.Migrations {
		var found *Migration
		for _, ref := range refsInDb {
			if migration.RefName == ref {
				found = migration
				break
			}
		}

		if found != nil {
			pendingMigrations = append(pendingMigrations, found)
		}
	}

	return pendingMigrations
}

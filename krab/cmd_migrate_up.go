package krab

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ohkrab/krab/krabdb"
	"github.com/ohkrab/krab/krabtpl"
	"github.com/ohkrab/krab/tpls"
	"github.com/pkg/errors"
)

// CmdMigrateUp returns migration status information.
type CmdMigrateUp struct {
	Set        *MigrationSet
	Connection krabdb.Connection
}

// ResponseMigrateUp json
type ResponseMigrateUp struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Success bool   `json:"success"`
}

func (c *CmdMigrateUp) Name() []string { return []string{"migrate", "up", c.Set.RefName} }

func (c *CmdMigrateUp) HttpMethod() string { return http.MethodPost }

func (c *CmdMigrateUp) Do(ctx context.Context, o CmdOpts) (interface{}, error) {
	for _, arg := range c.Set.Arguments.Args {
		_, ok := o.Inputs[arg.Name]
		if !ok {
			return nil, fmt.Errorf("Command is missing an input for argument `%s`", arg.Name)
		}
	}

	err := c.Set.Arguments.Validate(o.Inputs)
	if err != nil {
		return nil, err
	}

	var result []ResponseMigrateUp
	err = c.Connection.Get(func(db krabdb.DB) error {
		resp, err := c.run(ctx, db, o.Inputs)
		result = resp
		return err
	})

	return result, err
}

func (c *CmdMigrateUp) run(ctx context.Context, db krabdb.DB, inputs Inputs) ([]ResponseMigrateUp, error) {
	result := []ResponseMigrateUp{}

	tpl := tpls.New(inputs, krabtpl.Functions)
	versions := NewSchemaMigrationTable(tpl.Render(c.Set.Schema))

	// locking
	lockID := int64(1)

	_, err := krabdb.TryAdvisoryLock(ctx, db, lockID)
	if err != nil {
		return nil, errors.Wrap(err, "Possibly another migration in progress")
	}
	defer krabdb.AdvisoryUnlock(ctx, db, lockID)

	hooksRunner := HookRunner{}
	err = hooksRunner.SetSearchPath(ctx, db, tpl.Render(c.Set.Schema))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to run SetSearchPath hook")
	}

	// schema migration
	err = versions.Init(ctx, db)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create default table for migrations")
	}

	migrationRefsInDb, err := versions.SelectAll(ctx, db)
	if err != nil {
		return nil, err
	}

	pendingMigrations := versions.FilterPending(c.Set.Migrations, migrationRefsInDb)

	for _, pending := range pendingMigrations {
		tx, err := db.NewTx(ctx, pending.ShouldRunInTransaction())
		if err != nil {
			result = append(result, ResponseMigrateUp{
				Name:    pending.RefName,
				Version: pending.Version,
				Success: false,
			})
			return result, errors.Wrap(err, "Failed to start transaction")
		}
		err = hooksRunner.SetSearchPath(ctx, tx, tpl.Render(c.Set.Schema))
		if err != nil {
			result = append(result, ResponseMigrateUp{
				Name:    pending.RefName,
				Version: pending.Version,
				Success: false,
			})
			return result, errors.Wrap(err, "Failed to run SetSearchPath hook")
		}

		err = c.migrateUp(ctx, tx, pending, versions)
		if err != nil {
			result = append(result, ResponseMigrateUp{
				Name:    pending.RefName,
				Version: pending.Version,
				Success: false,
			})
			tx.Rollback()
			return result, err
		}

		err = tx.Commit()
		if err != nil {
			result = append(result, ResponseMigrateUp{
				Name:    pending.RefName,
				Version: pending.Version,
				Success: false,
			})
			return result, err
		}

		result = append(result, ResponseMigrateUp{
			Name:    pending.RefName,
			Version: pending.Version,
			Success: true,
		})
	}

	return result, nil
}

func (c *CmdMigrateUp) migrateUp(ctx context.Context, tx krabdb.TransactionExecerContext, migration *Migration, versions SchemaMigrationTable) error {
	sqls := migration.Up.ToSQLStatements()
	for _, sql := range sqls {
		// fmt.Println(ctc.ForegroundYellow, string(sql), ctc.Reset)
		_, err := tx.ExecContext(ctx, string(sql))
		if err != nil {
			return errors.Wrap(err, "Failed to execute migration")
		}
	}

	err := versions.Insert(ctx, tx, migration.Version)
	if err != nil {
		return errors.Wrap(err, "Failed to insert migration")
	}

	return nil
}

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

// CmdMigrateDown returns migration status information.
type CmdMigrateDown struct {
	Set        *MigrationSet
	Connection krabdb.Connection
}

// ResponseMigrateDown json
type ResponseMigrateDown struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Success bool   `json:"success"`
}

func (c *CmdMigrateDown) Arguments() *Arguments {
	return &Arguments{
		Args: []*Argument{
			{
				Name:        "version",
				Type:        "string",
				Description: "Migration version to rollback",
			},
		},
	}
}

func (c *CmdMigrateDown) Name() []string { return []string{"migrate", "down", c.Set.RefName} }

func (c *CmdMigrateDown) HttpMethod() string { return http.MethodPost }

func (c *CmdMigrateDown) Do(ctx context.Context, o CmdOpts) (interface{}, error) {
	for _, arg := range c.Set.Arguments.Args {
		_, ok := o.Inputs[arg.Name]
		if !ok {
			return nil, fmt.Errorf("Command is missing an input for argument `%s`", arg.Name)
		}
	}
	// default arguments always take precedence over custom ones
	for _, arg := range c.Arguments().Args {
		_, ok := o.Inputs[arg.Name]
		if !ok {
			return nil, fmt.Errorf("Command is missing an input for argument `%s`", arg.Name)
		}
	}

	err := c.Set.Arguments.Validate(o.Inputs)
	if err != nil {
		return nil, err
	}
	err = c.Arguments().Validate(o.Inputs)
	if err != nil {
		return nil, err
	}

	var result []ResponseMigrateDown
	err = c.Connection.Get(func(db krabdb.DB) error {
		resp, err := c.run(ctx, db, o.Inputs)
		result = resp
		return err
	})

	return result, err
}

func (c *CmdMigrateDown) run(ctx context.Context, db krabdb.DB, inputs Inputs) ([]ResponseMigrateDown, error) {
	result := []ResponseMigrateDown{}

	tpl := tpls.New(inputs, krabtpl.Functions)
	versions := NewSchemaMigrationTable(tpl.Render(c.Set.Schema))

	migration := c.Set.FindMigrationByVersion(inputs["version"].(string))
	if migration == nil {
		return nil, fmt.Errorf("Migration `%s` not found in `%s` set",
			inputs["version"].(string),
			c.Set.RefName)
	}
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
	tx, err := db.NewTx(ctx, migration.ShouldRunInTransaction())
	if err != nil {
		return nil, errors.Wrap(err, "Failed to start transaction")
	}
	err = hooksRunner.SetSearchPath(ctx, tx, tpl.Render(c.Set.Schema))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to run SetSearchPath hook")
	}

	migrationExists, _ := versions.Exists(ctx, db, SchemaMigration{migration.Version})
	if migrationExists {
		sqls := migration.Down.ToSQLStatements()
		for _, sql := range sqls {
			// fmt.Println(ctc.ForegroundYellow, string(sql), ctc.Reset)
			_, err := tx.ExecContext(ctx, string(sql))
			if err != nil {
				result = append(result, ResponseMigrateDown{
					Name:    migration.RefName,
					Version: migration.Version,
					Success: false,
				})
				tx.Rollback()
				return nil, errors.Wrap(err, "Failed to execute migration")
			}
		}

		err := versions.Delete(ctx, tx, migration.Version)
		if err != nil {
			result = append(result, ResponseMigrateDown{
				Name:    migration.RefName,
				Version: migration.Version,
				Success: false,
			})
			tx.Rollback()
			return nil, errors.Wrap(err, "Failed to delete migration")
		}

		result = append(result, ResponseMigrateDown{
			Name:    migration.RefName,
			Version: migration.Version,
			Success: true,
		})
	} else {
		tx.Rollback()
		return nil, errors.New("Migration has not been run yet, nothing to rollback")
	}

	err = tx.Commit()

	return result, err
}

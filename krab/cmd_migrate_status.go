package krab

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/ohkrab/krab/krabdb"
	"github.com/ohkrab/krab/krabtpl"
	"github.com/ohkrab/krab/tpls"
	"github.com/pkg/errors"
)

// CmdMigrateStatus returns migration status information.
type CmdMigrateStatus struct {
	Set        *MigrationSet
	Connection krabdb.Connection
	Inputs
}

// ResponseMigrateStatus json
type ResponseMigrateStatus struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Pending bool   `json:"pending"`
}

func (c *CmdMigrateStatus) Name() []string { return []string{"migrate", "status", c.Set.RefName} }

func (c *CmdMigrateStatus) HttpMethod() string { return http.MethodGet }

func (c *CmdMigrateStatus) Do(ctx context.Context, o CmdOpts) error {
	for _, arg := range c.Set.Arguments.Args {
		_, ok := c.Inputs[arg.Name]
		if !ok {
			return fmt.Errorf("Command is missing an input for argument `%s`", arg.Name)
		}
	}
	err := c.Set.Arguments.Validate(c.Inputs)
	if err != nil {
		return err
	}

	err = c.Connection.Get(func(db krabdb.DB) error {
		resp, err := c.run(ctx, db)
		if err == nil {
			return json.NewEncoder(o.Writer).Encode(resp)
		}
		return err
	})

	return err
}

func (c *CmdMigrateStatus) run(ctx context.Context, db krabdb.DB) ([]ResponseMigrateStatus, error) {
	result := []ResponseMigrateStatus{}

	tpl := tpls.New(c.Inputs, krabtpl.Functions)
	versions := NewSchemaMigrationTable(tpl.Render(c.Set.Schema))

	hooksRunner := HookRunner{}
	err := hooksRunner.SetSearchPath(ctx, db, tpl.Render(c.Set.Schema))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to run SetSearchPath hook")
	}
	migrationRefsInDb, err := versions.SelectAll(ctx, db)
	if err != nil {
		return nil, err
	}

	appliedMigrations := hashset.New()

	for _, migration := range migrationRefsInDb {
		appliedMigrations.Add(migration.Version)
	}

	for _, migration := range c.Set.Migrations {
		pending := !appliedMigrations.Contains(migration.Version)

		if pending {
			result = append(result, ResponseMigrateStatus{
				Name:    migration.RefName,
				Version: migration.Version,
				Pending: true,
			})
		} else {
			result = append(result, ResponseMigrateStatus{
				Name:    migration.RefName,
				Version: migration.Version,
				Pending: false,
			})
		}

	}

	return result, nil
}

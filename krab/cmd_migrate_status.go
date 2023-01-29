package krab

import (
	"context"
	"net/http"

	"github.com/emirpasic/gods/sets/hashset"
	"github.com/ohkrab/krab/krabdb"
	"github.com/ohkrab/krab/krabhcl"
	"github.com/ohkrab/krab/krabtpl"
	"github.com/ohkrab/krab/tpls"
	"github.com/pkg/errors"
)

// CmdMigrateStatus returns migration status information.
type CmdMigrateStatus struct {
	Set        *MigrationSet
	Connection krabdb.Connection
}

// ResponseMigrateStatus json
type ResponseMigrateStatus struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Pending bool   `json:"pending"`
}

func (c *CmdMigrateStatus) Addr() krabhcl.Addr { return c.Set.Addr() }

func (c *CmdMigrateStatus) Name() []string {
	return append([]string{"migrate", "status"}, c.Set.Addr().Labels...)
}

func (c *CmdMigrateStatus) HttpMethod() string { return http.MethodGet }

func (c *CmdMigrateStatus) Do(ctx context.Context, o CmdOpts) (interface{}, error) {
	err := c.Set.Arguments.Validate(o.NamedInputs)
	if err != nil {
		return nil, err
	}

	var result []ResponseMigrateStatus
	err = c.Connection.Get(func(db krabdb.DB) error {
		resp, err := c.run(ctx, db, o.NamedInputs)
		result = resp
		return err
	})

	return result, err
}

func (c *CmdMigrateStatus) run(ctx context.Context, db krabdb.DB, inputs NamedInputs) ([]ResponseMigrateStatus, error) {
	result := []ResponseMigrateStatus{}

	tpl := tpls.New(inputs, krabtpl.Functions)
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

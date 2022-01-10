package krab

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ohkrab/krab/krabdb"
	"github.com/ohkrab/krab/krabhcl"
	"github.com/wzshiming/ctc"
)

// CmdTestRun returns migration status information.
type CmdTestRun struct {
	Connection krabdb.Connection
	Suite      *TestSuite
	Registry   *CmdRegistry
}

// ResponseTestRun json
type ResponseTestRun struct {
}

func (c *CmdTestRun) Addr() krabhcl.Addr { return c.Suite.Addr() }

func (c *CmdTestRun) Name() []string { return append([]string{"test"}, c.Suite.Addr().Labels...) }

func (c *CmdTestRun) HttpMethod() string { return http.MethodPost }

func (c *CmdTestRun) Do(ctx context.Context, o CmdOpts) (interface{}, error) {
	var result ResponseTestRun

	for _, do := range c.Suite.Before.Dos {
		addr, err := krabhcl.ExpressionToAddr(do.MigrationSet)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse MigrationSet reference: %w", err)
		}

		for _, cmd := range c.Registry.Commands {
			if addr.Equal(cmd.Addr()) {
				if cmd.Name()[1] == "up" {
					result, err := cmd.Do(ctx, CmdOpts{Inputs: do.Inputs()})
					if err != nil {
						return nil, fmt.Errorf("Failed to execute before hook: %w", err)
					}
					resp := result.([]ResponseMigrateUp)
					for _, migration := range resp {
						fmt.Println(ctc.ForegroundYellow, migration.Success, migration.Version, migration.Name, ctc.Reset)
					}
				}
			}
		}
	}

	// err := c.Connection.Get(func(db krabdb.DB) error {
	// 	resp, err := c.run(ctx, db, o.Inputs)
	// 	result = resp
	// 	return err
	// })

	return result, nil
}

func (c *CmdTestRun) run(ctx context.Context, db krabdb.DB, inputs Inputs) (ResponseTestRun, error) {
	result := ResponseTestRun{}

	// _, err := db.ExecContext(ctx, sql)

	return result, nil
}

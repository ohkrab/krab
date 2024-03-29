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
}

// ResponseTestRun json
type ResponseTestRun struct {
}

func (c *CmdTestRun) Addr() krabhcl.Addr { return c.Suite.Addr() }

func (c *CmdTestRun) Name() []string { return []string{"test"} }

func (c *CmdTestRun) HttpMethod() string { return http.MethodPost }

func (c *CmdTestRun) Do(ctx context.Context, o CmdOpts) (interface{}, error) {
	var result ResponseTestRun

	// for _, do := range c.Suite.Before.Dos {
	// 	for _, migrate := range do.Migrate {
	// 		addr, err := krabhcl.Expression{Expr: migrate.SetExpr}.Addr()
	// 		if err != nil {
	// 			return nil, fmt.Errorf("Failed to parse MigrationSet reference: %w", err)
	// 		}
	//
	// 		for _, cmd := range c.Registry.Commands {
	// 			if addr.Equal(cmd.Addr()) {
	// 				if cmd.Name()[1] == migrate.Type {
	// 					inputs := InputsFromCtyInputs(do.CtyInputs)
	// 					migrateInputs := InputsFromCtyInputs(migrate.CtyInputs)
	// 					inputs.Merge(migrateInputs)
	// 					result, err := cmd.Do(ctx, CmdOpts{NamedInputs: inputs})
	// 					if err != nil {
	// 						return nil, fmt.Errorf("Failed to execute before hook: %w", err)
	// 					}
	// 					respUp, ok := result.([]ResponseMigrateUp)
	// 					if ok {
	// 						for _, migration := range respUp {
	// 							fmt.Println(ctc.ForegroundYellow, "UP  ", migration.Success, migration.Version, migration.Name, ctc.Reset)
	// 						}
	// 					}
	// 					respDown, ok := result.([]ResponseMigrateDown)
	// 					if ok {
	// 						for _, migration := range respDown {
	// 							fmt.Println(ctc.ForegroundYellow, "DOWN", migration.Success, migration.Version, migration.Name, ctc.Reset)
	// 						}
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	err := c.Connection.Get(func(db krabdb.DB) error {
		resp, err := c.run(ctx, db, o.NamedInputs)
		result = resp
		return err
	})

	return result, err
}

func (c *CmdTestRun) run(ctx context.Context, db krabdb.DB, inputs NamedInputs) (ResponseTestRun, error) {
	result := ResponseTestRun{}

	for _, testCase := range c.Suite.Tests {
		fmt.Printf("%s%s%s\n", ctc.ForegroundBlue, testCase.Name, ctc.Reset)
		for _, it := range testCase.Xits {
			fmt.Printf("  %sSKIP %s%s\n", ctc.ForegroundYellow, it.Name, ctc.Reset)
		}
		for _, it := range testCase.Its {
			fmt.Printf("  %s%s%s\n", ctc.ForegroundRed, it.Name, ctc.Reset)
		}
	}

	// _, err := db.ExecContext(ctx, sql)

	return result, nil
}

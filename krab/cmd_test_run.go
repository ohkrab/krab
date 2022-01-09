package krab

import (
	"context"
	"net/http"

	"github.com/ohkrab/krab/krabdb"
)

// CmdTestRun returns migration status information.
type CmdTestRun struct {
	Connection krabdb.Connection
	Suite      *TestSuite
}

// ResponseTestRun json
type ResponseTestRun struct {
}

func (c *CmdTestRun) Name() []string { return []string{"test", c.Suite.RefName} }

func (c *CmdTestRun) HttpMethod() string { return http.MethodPost }

func (c *CmdTestRun) Do(ctx context.Context, o CmdOpts) (interface{}, error) {
	var result ResponseTestRun

	// for _, do := range c.Suite.Before.Dos {
	// }
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

package krab

import (
	"context"
	"net/http"
	"strings"

	"github.com/ohkrab/krab/krabdb"
	"github.com/ohkrab/krab/krabtpl"
	"github.com/ohkrab/krab/tpls"
)

// CmdAction returns migration status information.
type CmdAction struct {
	Action     *Action
	Connection krabdb.Connection
}

// ResponseAction json
type ResponseAction struct{}

func (c *CmdAction) Name() []string { return []string{"action", c.Action.Namespace, c.Action.RefName} }

func (c *CmdAction) HttpMethod() string { return http.MethodPost }

func (c *CmdAction) Do(ctx context.Context, o CmdOpts) (interface{}, error) {
	err := c.Action.Arguments.Validate(o.Inputs)
	if err != nil {
		return nil, err
	}

	var result ResponseAction
	err = c.Connection.Get(func(db krabdb.DB) error {
		resp, err := c.run(ctx, db, o.Inputs)
		result = resp
		return err
	})

	return result, err
}

func (c *CmdAction) run(ctx context.Context, db krabdb.DB, inputs Inputs) (ResponseAction, error) {
	result := ResponseAction{}

	tpl := tpls.New(inputs, krabtpl.Functions)

	sb := strings.Builder{}
	c.Action.ToSQL(&sb)
	sql := tpl.Render(sb.String())

	_, err := db.ExecContext(ctx, sql)

	return result, err
}

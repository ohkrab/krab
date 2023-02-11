package krab

import (
	"context"
	"net/http"
	"strings"

	"github.com/ohkrab/krab/krabdb"
	"github.com/ohkrab/krab/krabhcl"
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

func (c *CmdAction) Addr() krabhcl.Addr { return c.Action.Addr() }

func (c *CmdAction) Name() []string { return append([]string{"action"}, c.Action.Addr().Labels...) }

func (c *CmdAction) HttpMethod() string { return http.MethodPost }

func (c *CmdAction) Do(ctx context.Context, o CmdOpts) (interface{}, error) {
	err := c.Action.Arguments.Validate(o.NamedInputs)
	if err != nil {
		return nil, err
	}

	var result ResponseAction
	err = c.Connection.Get(func(db krabdb.DB) error {
		resp, err := c.run(ctx, db, o.NamedInputs)
		result = resp
		return err
	})

	return result, err
}

func (c *CmdAction) run(ctx context.Context, db krabdb.DB, inputs NamedInputs) (ResponseAction, error) {
	result := ResponseAction{}

	tpl := tpls.New(inputs, krabtpl.Functions())

	sb := strings.Builder{}
	c.Action.ToSQL(&sb)
	sql := tpl.Render(sb.String())

	_, err := db.ExecContext(ctx, sql)

	return result, err
}

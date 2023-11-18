package krab

import (
	"fmt"
	"io"

	"github.com/hashicorp/hcl/v2"
	"github.com/ohkrab/krab/krabhcl"
)

// Action represents custom action to execute.
type Action struct {
	krabhcl.Source

	Namespace string
	RefName   string

	Arguments *Arguments

	Description string
	SQL         string
	Transaction bool // wrap operation in transaction
}

func (a *Action) Addr() krabhcl.Addr {
	return krabhcl.Addr{Keyword: "action", Labels: []string{a.Namespace, a.RefName}}
}

var schemaAction = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "arguments",
			LabelNames: []string{},
		},
	},
	Attributes: []hcl.AttributeSchema{
		{
			Name:     "sql",
			Required: true,
		},
		{
			Name:     "transaction",
			Required: false,
		},
		{
			Name:     "description",
			Required: true,
		},
	},
}

// DecodeHCL parses HCL into struct.
func (a *Action) DecodeHCL(ctx *hcl.EvalContext, block *hcl.Block) error {
	a.Source.Extract(block)

	a.Namespace = block.Labels[0]
	a.RefName = block.Labels[1]
	a.Arguments = &Arguments{}
	a.Transaction = true

	content, diags := block.Body.Content(schemaAction)
	if diags.HasErrors() {
		return fmt.Errorf("failed to decode `%s` block: %s", block.Type, diags.Error())
	}

	for _, b := range content.Blocks {
		switch b.Type {
		case "arguments":
			err := a.Arguments.DecodeHCL(ctx, b)
			if err != nil {
				return err
			}

		default:
			return fmt.Errorf("Unknown block `%s` for `%s` block", b.Type, block.Type)
		}
	}

	for k, v := range content.Attributes {
		switch k {
		case "sql":
			expr := krabhcl.Expression{Expr: v.Expr, EvalContext: ctx}
			val, err := expr.String()
			if err != nil {
				return err
			}
			a.SQL = val

		case "description":
			expr := krabhcl.Expression{Expr: v.Expr, EvalContext: ctx}
			val, err := expr.String()
			if err != nil {
				return err
			}
			a.Description = val

		case "transaction":
			expr := krabhcl.Expression{Expr: v.Expr, EvalContext: ctx}
			val, err := expr.Bool()
			if err != nil {
				return err
			}
			a.Transaction = val

		default:
			return fmt.Errorf("Unknown attribute `%s` for `%s` block", k, block.Type)
		}
	}

	return nil
}

func (a *Action) Validate() error {
	return ErrorCoalesce(
		ValidateRefName(a.Namespace),
		ValidateRefName(a.RefName),
		ValidateStringNonEmpty("description", a.Description),
	)
}

func (m *Action) ToSQL(w io.StringWriter) {
	w.WriteString(m.SQL)
}

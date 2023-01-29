package krab

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/ohkrab/krab/krabhcl"
)

type Argument struct {
	Name        string
	Type        string
	Description string
}

var schemaArgument = hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{},
	Attributes: []hcl.AttributeSchema{
		{
			Name:     "type",
			Required: false,
		},
		{
			Name:     "description",
			Required: false,
		},
	},
}

// Arguments represents command line arguments or params that you can pass to action.
//
type Arguments struct {
	Args []*Argument
}

var schemaArguments = hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "arg",
			LabelNames: []string{"name"},
		},
	},
}

// DecodeHCL parses HCL into struct.
func (a *Arguments) DecodeHCL(ctx *hcl.EvalContext, block *hcl.Block) error {
	a.Args = []*Argument{}

	content, diags := block.Body.Content(&schemaArguments)
	if diags.HasErrors() {
		return fmt.Errorf("failed to decode `%s` block: %s", block.Type, diags.Error())
	}

	for _, b := range content.Blocks {
		switch b.Type {
		case "arg":
			arg := new(Argument)
			err := arg.DecodeHCL(ctx, b)
			if err != nil {
				return err
			}
			a.Args = append(a.Args, arg)

		default:
			return fmt.Errorf("Unknown block `%s` for `%s` block", b.Type, block.Type)
		}
	}

	return nil
}

func (a *Arguments) Validate(values NamedInputs) error {
	for _, a := range a.Args {
		value, ok := values[a.Name]
		if ok {
			if err := a.Validate(value); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("Argument value for `%s` (%s) is missing", a.Name, a.Description)
		}
	}

	return nil
}

func (a *Arguments) Help() string {
	sb := strings.Builder{}
	if len(a.Args) > 0 {
		for _, arg := range a.Args {
			sb.WriteString("  -")
			sb.WriteString(arg.Name)
			sb.WriteString(" (")
			sb.WriteString(arg.Type)
			sb.WriteString(", required")
			sb.WriteString(") ")
			sb.WriteString(arg.Description)
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// DecodeHCL parses HCL into struct.
func (a *Argument) DecodeHCL(ctx *hcl.EvalContext, block *hcl.Block) error {
	a.Name = block.Labels[0]
	a.Type = "string"
	a.Description = ""

	content, diags := block.Body.Content(&schemaArgument)
	if diags.HasErrors() {
		return fmt.Errorf("failed to decode `%s` block: %s", block.Type, diags.Error())
	}

	// no blocks to decode

	for k, v := range content.Attributes {
		switch k {
		case "type":
			expr := krabhcl.Expression{Expr: v.Expr, EvalContext: ctx}
			val, err := expr.String()
			if err != nil {
				return err
			}
			a.Type = val

		case "description":
			expr := krabhcl.Expression{Expr: v.Expr, EvalContext: ctx}
			val, err := expr.String()
			if err != nil {
				return err
			}
			a.Description = val

		default:
			return fmt.Errorf("Unknown attribute `%s` for `migration` block", k)
		}
	}

	return nil
}

func (a *Argument) Validate(value interface{}) error {
	switch value.(type) {
	case string:
		if len(value.(string)) == 0 {
			return fmt.Errorf("Value for -%s is required", a.Name)
		}
	default:
		return errors.New("Argument type not implemented")
	}
	return nil
}

package krab

import (
	"errors"
	"fmt"
	"strings"
)

type Argument struct {
	Name        string `hcl:"name,label"`
	Type        string `hcl:"type,optional"`
	Description string `hcl:"description,optional"`
}

// Arguments represents command line arguments or params that you can pass to action.
//
type Arguments struct {
	Args []*Argument `hcl:"arg,block"`
}

func (a *Arguments) Validate(values Inputs) error {
	for _, a := range a.Args {
		value, ok := values[a.Name]
		if ok {
			if err := a.Validate(value); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("Argument value for `%s` is missing", a.Name)
		}
	}

	return nil
}

func (a *Arguments) InitDefaults() {
	for _, a := range a.Args {
		a.InitDefaults()
	}
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

func (a *Argument) InitDefaults() {
	if a.Type == "" {
		a.Type = "string"
	}
}
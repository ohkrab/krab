package krab

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

// Do subtype for other types.
type Do struct {
	Migrate   []*DoMigrate         `hcl:"migrate,block"`
	CtyInputs map[string]cty.Value `hcl:"inputs,optional"`
	SQL       string               `hcl:"sql,optional"`
}

func (d *Do) Validate() error {
	for _, m := range d.Migrate {
		if err := m.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type DoMigrate struct {
	Type      string               `hcl:"type,label"`
	SetExpr   hcl.Expression       `hcl:"migration_set"`
	CtyInputs map[string]cty.Value `hcl:"inputs,optional"`

	Set *MigrationSet
}

func (d *DoMigrate) Validate() error {
	switch d.Type {
	case "up", "down":
		return nil
	}

	return fmt.Errorf("Invalid type `%s` for `do` command", d.Type)
}

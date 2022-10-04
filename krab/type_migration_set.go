package krab

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/ohkrab/krab/krabhcl"
)

// MigrationSet represents collection of migrations.
type MigrationSet struct {
	krabhcl.Source

	RefName string
	Schema  string
	// SchemaMigrationTableName `hcl:"schema_migrations_table,optional"`

	Arguments *Arguments
	Hooks     *Hooks

	MigrationAddrs []*krabhcl.Addr
	Migrations     []*Migration // populated from refs in expression
}

var schemaMigrationSet = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "arguments",
			LabelNames: []string{},
		},
	},
	Attributes: []hcl.AttributeSchema{
		{
			Name:     "migrations",
			Required: true,
		},
		{
			Name:     "schema",
			Required: false,
		},
	},
}

func (ms *MigrationSet) Addr() krabhcl.Addr {
	return krabhcl.Addr{Keyword: "migration_set", Labels: []string{ms.RefName}}
}

// DecodeHCL parses HCL into struct.
func (ms *MigrationSet) DecodeHCL(ctx *hcl.EvalContext, block *hcl.Block) error {
	ms.Source.Extract(block)

	ms.RefName = block.Labels[0]
	ms.Schema = "public"
	ms.Arguments = &Arguments{}
	ms.Hooks = &Hooks{}
	ms.Migrations = []*Migration{}
	ms.MigrationAddrs = []*krabhcl.Addr{}

	content, diags := block.Body.Content(schemaMigrationSet)
	if diags.HasErrors() {
		return fmt.Errorf("failed to decode `%s` block: %s", block.Type, diags.Error())
	}

	for _, b := range content.Blocks {
		switch b.Type {
		case "arguments":
			err := ms.Arguments.DecodeHCL(ctx, b)
			if err != nil {
				return err
			}

		default:
			return fmt.Errorf("Unknown block `%s` for `%s` block", b.Type, block.Type)
		}
	}

	for k, v := range content.Attributes {
		switch k {
		case "migrations":
			expr := krabhcl.Expression{Expr: v.Expr, EvalContext: ctx}
			val, err := expr.SliceAddr()
			if err != nil {
				return err
			}
			ms.MigrationAddrs = append(ms.MigrationAddrs, val...)

		case "schema":
			expr := krabhcl.Expression{Expr: v.Expr, EvalContext: ctx}
			val, err := expr.String()
			if err != nil {
				return err
			}
			ms.Schema = val

		default:
			return fmt.Errorf("Unknown attribute `%s` for `migration_set` block", k)
		}
	}

	return nil
}

func (ms *MigrationSet) Validate() error {
	return ErrorCoalesce(
		ValidateRefName(ms.RefName),
	)
}

// FindMigrationByVersion looks up for the migration in current set.
func (ms *MigrationSet) FindMigrationByVersion(version string) *Migration {
	for _, m := range ms.Migrations {
		if m.Version == version {
			return m
		}
	}

	return nil
}

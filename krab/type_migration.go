package krab

import (
	"fmt"
	"io"

	"github.com/hashicorp/hcl/v2"
	"github.com/ohkrab/krab/krabhcl"
)

// Migration represents single up/down migration pair.
//
type Migration struct {
	krabhcl.Source

	RefName     string
	Version     string
	Up          MigrationUpOrDown
	Down        MigrationUpOrDown
	Transaction bool // wrap operation in transaction
}

var schemaMigration = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "up",
			LabelNames: []string{},
		},
		{
			Type:       "down",
			LabelNames: []string{},
		},
	},
	Attributes: []hcl.AttributeSchema{
		{
			Name:     "version",
			Required: true,
		},
		{
			Name:     "transaction",
			Required: false,
		},
	},
}

// Migration contains info how to migrate up or down.
type MigrationUpOrDown struct {
	krabhcl.Source

	SQL           string            `hcl:"sql,optional"`
	CreateTables  []*DDLCreateTable `hcl:"create_table,block"`
	CreateIndices []*DDLCreateIndex `hcl:"create_index,block"`
	DropTables    []*DDLDropTable   `hcl:"drop_table,block"`
	DropIndices   []*DDLDropIndex   `hcl:"drop_index,block"`

	AttrDefRanges map[string]hcl.Range
}

var schemaMigrationUpOrDown = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "create_table",
			LabelNames: []string{"name"},
		},
		{
			Type:       "create_index",
			LabelNames: []string{"table", "name"},
		},
		{
			Type:       "drop_table",
			LabelNames: []string{"name"},
		},
		{
			Type:       "drop_index",
			LabelNames: []string{"name"},
		},
	},
	Attributes: []hcl.AttributeSchema{
		{
			Name:     "sql",
			Required: false,
		},
	},
}

// DecodeHCL parses HCL into struct.
func (m *Migration) DecodeHCL(ctx *hcl.EvalContext, block *hcl.Block) error {
	m.Source.Extract(block)

	m.RefName = block.Labels[0]

	// set defaults and init
	m.Transaction = true

	content, diags := block.Body.Content(schemaMigration)
	if diags.HasErrors() {
		return fmt.Errorf("failed to decode `migration` block: %s", diags.Error())
	}

	for _, b := range content.Blocks {
		switch b.Type {
		case "up":
			err := m.Up.DecodeHCL(ctx, b)
			if err != nil {
				return err
			}
		case "down":
			err := m.Down.DecodeHCL(ctx, b)
			if err != nil {
				return err
			}
		}
	}

	for k, v := range content.Attributes {
		switch k {
		case "version":
			expr := krabhcl.Expression{Expr: v.Expr, EvalContext: ctx}
			m.Version = expr.AsString()

		case "transaction":
			expr := krabhcl.Expression{Expr: v.Expr, EvalContext: ctx}
			m.Transaction = expr.AsBool()

		default:
			return fmt.Errorf("Unknown attribute `%s` for `migration` block", k)
		}
	}

	return nil
}

func (ms *Migration) Validate() error {
	return ErrorCoalesce(
		ValidateRefName(ms.RefName),
		ms.Up.Validate(),
		ms.Down.Validate(),
	)
}

// DecodeHCL parses HCL into struct.
func (m *MigrationUpOrDown) DecodeHCL(ctx *hcl.EvalContext, block *hcl.Block) error {
	m.Source.Extract(block)
	m.AttrDefRanges = map[string]hcl.Range{}

	content, diags := block.Body.Content(schemaMigrationUpOrDown)
	if diags.HasErrors() {
		return fmt.Errorf("failed to decode `%s` block: %s", block.Type, diags.Error())
	}

	for _, b := range content.Blocks {
		switch b.Type {
		case "create_table":
		case "create_index":
		case "drop_table":
		case "drop_index":
		}
	}

	for k, v := range content.Attributes {
		switch k {
		case "sql":
			m.AttrDefRanges["sql"] = v.Range
			expr := krabhcl.Expression{Expr: v.Expr, EvalContext: ctx}
			m.SQL = expr.AsString()

		default:
			return fmt.Errorf("Unknown attribute `%s` for `%s` block", k, block.Type)
		}
	}

	return nil
}

func (m *MigrationUpOrDown) Validate() error {
	return nil
}

func (m *MigrationUpOrDown) ToSQL(w io.StringWriter) {
	w.WriteString(m.SQL)
}

// ToSQLStatements returns list of SQL statements to executre during the migration.
func (m *MigrationUpOrDown) ToSQLStatements() SQLStatements {
	sorter := SQLStatementsSorter{Statements: SQLStatements{}, Bytes: []int{}}

	if m.SQL != "" {
		sorter.Insert(m.AttrDefRanges["sql"], m)
	}

	for _, t := range m.CreateTables {
		sorter.Insert(t.DefRange, t)
	}
	for _, t := range m.CreateIndices {
		sorter.Insert(t.DefRange, t)
	}
	for _, t := range m.DropIndices {
		sorter.Insert(t.DefRange, t)
	}
	for _, t := range m.DropTables {
		sorter.Insert(t.DefRange, t)
	}

	return sorter.Sort()
}

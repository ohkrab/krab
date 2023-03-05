package krab

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/ohkrab/krab/krabhcl"
)

// TestExample represents test runner configuration.
type TestExample struct {
	krabhcl.Source

	// Set              *SetRuntimeParameters `hcl:"set,block"`
	Name string
	Its  []*TestExampleIt
	Xits []*TestExampleIt
}

var schemaTestExample = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "it",
			LabelNames: []string{"name"},
		},
		{
			Type:       "xit",
			LabelNames: []string{"name"},
		},
	},
	Attributes: []hcl.AttributeSchema{},
}

func (t *TestExample) Addr() krabhcl.Addr {
	return krabhcl.Addr{Keyword: "test", Labels: []string{t.Name}}
}

func (t *TestExample) Validate() error {
	return ErrorCoalesce()
}

func (t *TestExample) DecodeHCL(ctx *hcl.EvalContext, block *hcl.Block) error {
	t.Source.Extract(block)
	t.Name = block.Labels[0]

	content, diags := block.Body.Content(schemaTestExample)
	if diags.HasErrors() {
		return fmt.Errorf("failed to decode `test` block: %s", diags.Error())
	}

	for _, b := range content.Blocks {
		switch b.Type {
		case "it":
			it := new(TestExampleIt)
			err := it.DecodeHCL(ctx, b)
			if err != nil {
				return err
			}
			t.Its = append(t.Its, it)

		case "xit":
			it := new(TestExampleIt)
			err := it.DecodeHCL(ctx, b)
			if err != nil {
				return err
			}
			t.Xits = append(t.Xits, it)

		default:
			return fmt.Errorf("Unknown block `%s` for `%s` block", b.Type, block.Type)
		}
	}

	for k := range content.Attributes {
		switch k {

		default:
			return fmt.Errorf("Unknown attribute `%s` for `migration` block", k)
		}
	}

	return nil
}

// TestExampleIt represents one use case for test example that contain queries and assertions.
type TestExampleIt struct {
	krabhcl.Source

	Name    string
	Queries []*TestQuery
}

var schemaTestExampleIt = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "query",
			LabelNames: []string{"sql"},
		},
	},
	Attributes: []hcl.AttributeSchema{},
}

func (it *TestExampleIt) DecodeHCL(ctx *hcl.EvalContext, block *hcl.Block) error {
	it.Source.Extract(block)

	it.Name = block.Labels[0]

	content, diags := block.Body.Content(schemaTestExampleIt)
	if diags.HasErrors() {
		return fmt.Errorf("failed to decode `test` block: %s", diags.Error())
	}

	for _, b := range content.Blocks {
		switch b.Type {
		case "query":
			q := new(TestQuery)
			err := q.DecodeHCL(ctx, b)
			if err != nil {
				return err
			}
			it.Queries = append(it.Queries, q)

		default:
			return fmt.Errorf("Unknown block `%s` for `%s` block", b.Type, block.Type)
		}
	}

	for k, _ := range content.Attributes {
		switch k {

		default:
			return fmt.Errorf("Unknown attribute `%s` for `migration` block", k)
		}
	}

	return nil
}

type TestQuery struct {
	krabhcl.Source

	Query string
	Rows  []*TestQueryRow
}

var schemaTestQuery = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "row",
			LabelNames: []string{"scope"},
		},
	},
	Attributes: []hcl.AttributeSchema{},
}

func (q *TestQuery) DecodeHCL(ctx *hcl.EvalContext, block *hcl.Block) error {
	q.Source.Extract(block)
	q.Query = block.Labels[0]

	content, diags := block.Body.Content(schemaTestQuery)
	if diags.HasErrors() {
		return fmt.Errorf("failed to decode `test` block: %s", diags.Error())
	}

	for _, b := range content.Blocks {
		switch b.Type {
		case "row":
			row := new(TestQueryRow)
			err := row.DecodeHCL(ctx, b)
			if err != nil {
				return err
			}
			q.Rows = append(q.Rows, row)

		default:
			return fmt.Errorf("Unknown block `%s` for `%s` block", b.Type, block.Type)
		}
	}

	for k, _ := range content.Attributes {
		switch k {

		default:
			return fmt.Errorf("Unknown attribute `%s` for `migration` block", k)
		}
	}

	return nil
}

type TestQueryRow struct {
	krabhcl.Source

	Scope string
	Cols  []*TestQueryCol
}

var schemaTestQueryRow = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "col",
			LabelNames: []string{"message"},
		},
	},
	Attributes: []hcl.AttributeSchema{},
}

func (row *TestQueryRow) DecodeHCL(ctx *hcl.EvalContext, block *hcl.Block) error {
	row.Source.Extract(block)
	row.Scope = block.Labels[0]

	content, diags := block.Body.Content(schemaTestQueryRow)
	if diags.HasErrors() {
		return fmt.Errorf("failed to decode `test` block: %s", diags.Error())
	}

	for _, b := range content.Blocks {
		switch b.Type {
		case "col":
			col := new(TestQueryCol)
			err := col.DecodeHCL(ctx, b)
			if err != nil {
				return err
			}
			row.Cols = append(row.Cols, col)

		default:
			return fmt.Errorf("Unknown block `%s` for `%s` block", b.Type, block.Type)
		}
	}

	for k, _ := range content.Attributes {
		switch k {

		default:
			return fmt.Errorf("Unknown attribute `%s` for `migration` block", k)
		}
	}

	return nil
}

type TestQueryCol struct {
	krabhcl.Source

	Message string
	Assert  string
}

var schemaTestQueryCol = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
	},
	Attributes: []hcl.AttributeSchema{
		{
			Name:     "assert",
			Required: true,
		},
	},
}

func (col *TestQueryCol) DecodeHCL(ctx *hcl.EvalContext, block *hcl.Block) error {
	col.Source.Extract(block)
	col.Message = block.Labels[0]

	content, diags := block.Body.Content(schemaTestQueryCol)
	if diags.HasErrors() {
		return fmt.Errorf("failed to decode `test` block: %s", diags.Error())
	}

	for _, b := range content.Blocks {
		switch b.Type {

		default:
			return fmt.Errorf("Unknown block `%s` for `%s` block", b.Type, block.Type)
		}
	}

	for k, v := range content.Attributes {
		switch k {

		case "assert":
			expr := krabhcl.Expression{Expr: v.Expr, EvalContext: ctx}
			val, err := expr.String()
			if err != nil {
				return err
			}
			col.Assert = val

		default:
			return fmt.Errorf("Unknown attribute `%s` for `migration` block", k)
		}
	}

	return nil
}

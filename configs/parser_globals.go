package configs

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/ohkrab/krab/addrs"
)

type Global struct {
	addrs.Addr

	Name string
	Expr hcl.Expression

	DeclRange hcl.Range
}

var globalsBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{},
	Blocks:     []hcl.BlockHeaderSchema{},
}

// https://github.com/hashicorp/terraform/blob/e320cd2b357479e4a89cb9236ebbb4ca70e24dfc/configs/named_values.go#l492
func decodeGlobalsBlock(block *hcl.Block) ([]*Global, hcl.Diagnostics) {
	attrs, diags := block.Body.JustAttributes()
	if len(attrs) == 0 {
		return nil, diags
	}

	globals := make([]*Global, 0, len(attrs))
	for name, attr := range attrs {
		if !hclsyntax.ValidIdentifier(name) {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid global value name",
				Detail:   "TODO: invalid global format message",
				Subject:  &attr.NameRange,
			})
		}

		globals = append(globals, &Global{
			Addr: addrs.Addr{
				Keyword: "global",
				Type:    "",
				Name:    name,
			},
			Name:      name,
			Expr:      attr.Expr,
			DeclRange: attr.Range,
		})
	}

	return globals, diags
}

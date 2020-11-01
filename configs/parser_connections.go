package configs

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/ohkrab/krab/addrs"
)

type Connection struct {
	addrs.Addr

	Uri       hcl.Expression
	Config    hcl.Body
	DeclRange hcl.Range
	TypeRange hcl.Range
	UriVal    string
}

var connectionBlockSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{
			Name: "uri",
		},
	},
	Blocks: []hcl.BlockHeaderSchema{},
}

// https://github.com/hashicorp/terraform/blob/e320cd2b357479e4a89cb9236ebbb4ca70e24dfc/configs/named_values.go#l492
func decodeConnectionBlock(block *hcl.Block) (*Connection, hcl.Diagnostics) {
	c := &Connection{
		Addr: addrs.Addr{
			Keyword: "connection",
			Type:    "",
			Name:    block.Labels[0],
		},
	}
	content, remain, diags := block.Body.PartialContent(connectionBlockSchema)
	c.Config = remain

	if !hclsyntax.ValidIdentifier(c.Addr.Name) {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid identifier for connection name",
			Subject:  &block.LabelRanges[0],
		})
	}

	// param := function.Parameter{Name: "name", Type: cty.String}
	// fn := function.New(
	// 	&function.Spec{
	// 		Params: []function.Parameter{param},
	// 		Type:   function.StaticReturnType(cty.String),
	// 		Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
	// 			return cty.StringVal(os.Getenv(args[0].AsString())), nil
	// 		},
	// 	},
	// )

	if attr, exists := content.Attributes["uri"]; exists {
		c.Uri = attr.Expr
	}
	// evalContext := hcl.EvalContext{
	// 	Variables: map[string]cty.Value{
	// 		"local": cty.MapVal(map[string]cty.Value{"uri": cty.StringVal("postgres://uri")}),
	// 	},
	// 	Functions: map[string]function.Function{"env": fn},
	// }
	// tt := c.Uri.Variables()
	// for _, t := range tt {
	// 	fmt.Println("t", t)
	// 	a, b := t.TraverseAbs(&evalContext)
	// 	fmt.Println("T", t, a, b)
	// }
	// // if c.Name == "from_env" {
	// val, _ := c.Uri.Value(&evalContext)
	// c.UriVal = val.AsString()
	// }

	return c, diags
}

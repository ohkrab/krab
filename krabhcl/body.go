package krabhcl

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
)

type Body struct {
	hcl.Body
}

func (b *Body) DefRangesFromPartialContentBlocks(schema *hcl.BodySchema) []hcl.Range {
	ret := []hcl.Range{}
	content, _, diags := b.PartialContent(schema)
	if len(diags) > 0 {
		panic(fmt.Sprintf("krabhcl.Body failed to extract PartialContent from schema: %v", diags))
	}
	for _, block := range content.Blocks {
		ret = append(ret, block.DefRange)
	}

	return ret
}

func (b *Body) DefRangesFromPartialContentAttributes(schema *hcl.BodySchema) map[string]hcl.Range {
	ret := map[string]hcl.Range{}
	content, _, diags := b.PartialContent(schema)
	if len(diags) > 0 {
		panic(fmt.Sprintf("krabhcl.Body failed to extract PartialContent from schema: %v", diags))
	}
	for _, attr := range content.Attributes {
		ret[attr.Name] = attr.Range
	}

	return ret
}

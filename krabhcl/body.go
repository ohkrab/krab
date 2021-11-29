package krabhcl

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
)

type Body struct {
	hcl.Body
}

func (b *Body) DefRangesFromPartialContent(schema *hcl.BodySchema) []hcl.Range {
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

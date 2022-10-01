package krabhcl

import (
	"github.com/hashicorp/hcl/v2"
)

// Source helps identifing code definition in krah hcl files.
type Source struct {
	DefRange hcl.Range
}

// Extract saves the source information.
func (s *Source) Extract(block *hcl.Block) {
	s.DefRange = block.DefRange
}

package diagnostics

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
)

type SourceRange struct {
	Filename   string
	Start, End SourcePos
}

type SourcePos struct {
	Line, Column, Byte int
}

func (r SourceRange) String() string {
	filename := r.Filename

	// try extracting relative path
	wd, err := os.Getwd()
	if err == nil {
		relFn, err := filepath.Rel(wd, filename)
		if err == nil {
			filename = relFn
		}
	}

	return fmt.Sprintf("%s:%d,%d", filename, r.Start.Line, r.Start.Column)
}

func SourceRangeFromHCL(hclRange hcl.Range) SourceRange {
	return SourceRange{
		Filename: hclRange.Filename,
		Start: SourcePos{
			Line:   hclRange.Start.Line,
			Column: hclRange.Start.Column,
			Byte:   hclRange.Start.Byte,
		},
		End: SourcePos{
			Line:   hclRange.End.Line,
			Column: hclRange.End.Column,
			Byte:   hclRange.End.Byte,
		},
	}
}

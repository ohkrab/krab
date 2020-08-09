package diagnostics

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
)

type List []Diagnostic

type Diagnostic struct {
}

type Hcl struct {
	wrappedDiags *hcl.Diagnostic
}

func New() List {
	return List{}
}

func (list List) Append(diags ...interface{}) List {
	for _, item := range diags {
		if item == nil {
			continue
		}

		switch ti := item.(type) {
		case Diagnostic:
			list = append(list, ti)
		case List:
			list = append(list, ti...)
		case hcl.Diagnostics:
			for _, hclDiag := range ti {
				diags = append(diags, Hcl{hclDiag})
			}
		// case *hcl.Diagnostic:
		// 	diags = append(diags, hclDiagnostic{i})
		case error:
			switch {
			default:
				panic(fmt.Sprint("Error", ti))
				// diags = append(diags, ti)
			}
		default:
			panic("Diagnostic not implemented")
		}
	}

	return list
}

func (list List) HasErrors() bool {
	return len(list) > 0
}

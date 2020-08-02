package lang

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/ohkrab/krab/addrs"
	"github.com/ohkrab/krab/diagnostics"
)

func References(traversals []hcl.Traversal) ([]*addrs.Reference, diagnostics.List) {
	if len(traversals) == 0 {
		return nil, nil
	}

	var diags diagnostics.List
	refs := make([]*addrs.Reference, 0, len(traversals))

	for _, traversal := range traversals {
		ref, refDiags := addrs.ParseRef(traversal)
		diags = diags.Append(refDiags)
		if ref == nil {
			continue
		}
		refs = append(refs, ref)
	}

	return refs, diags
}

func ReferencesInExpr(expr hcl.Expression) ([]*addrs.Reference, diagnostics.List) {
	if expr == nil {
		return nil, nil
	}
	traversals := expr.Variables()
	return References(traversals)
}

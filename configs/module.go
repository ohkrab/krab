package configs

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

type Module struct {
	SourceDir   string
	Connections map[string]*Connection
}

type File struct {
	Connections []*Connection

	Variables map[string]cty.Value
	Functions map[string]function.Function
}

func NewModule(files []*File) (*Module, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	mod := &Module{
		Connections: map[string]*Connection{},
	}

	for _, file := range files {
		fileDiags := mod.appendFile(file)
		diags = append(diags, fileDiags...)
	}

	return mod, diags
}

func (m *Module) appendFile(file *File) hcl.Diagnostics {
	var diags hcl.Diagnostics

	// for _, l := range file.Locals {
	// 	if existing, exists := m.Locals[l.Name]; exists {
	// 		diags = append(diags, &hcl.Diagnostic{
	// 			Severity: hcl.DiagError,
	// 			Summary:  "Duplicate local value definition",
	// 			Detail:   fmt.Sprintf("A local value named %q was already defined at %s. Local value names must be unique within a module.", existing.Name, existing.DeclRange),
	// 			Subject:  &l.DeclRange,
	// 		})
	// 	}
	// 	m.Locals[l.Name] = l
	// }

	for _, r := range file.Connections {
		key := r.Addr.String()
		if existing, exists := m.Connections[key]; exists {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("Duplicate resource %q configuration", existing.Type),
				Detail:   fmt.Sprintf("A %s resource named %q was already declared at %s. Resource names must be unique per type in each module.", existing.Type, existing.Name, existing.DeclRange),
				Subject:  &r.DeclRange,
			})
			continue
		}
		m.Connections[key] = r
	}

	return diags
}

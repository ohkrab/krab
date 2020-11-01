package configs

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
)

type Module struct {
	SourceDir   string
	Connections map[string]*Connection
	Globals     map[string]*Global
}

type File struct {
	Connections []*Connection
	Globals     []*Global
}

func NewModule(files []*File) (*Module, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	mod := &Module{
		Connections: map[string]*Connection{},
		Globals:     map[string]*Global{},
	}

	for _, file := range files {
		fileDiags := mod.appendFile(file)
		diags = append(diags, fileDiags...)
	}

	return mod, diags
}

func (m *Module) appendFile(file *File) hcl.Diagnostics {
	var diags hcl.Diagnostics

	for _, g := range file.Globals {
		if existing, exists := m.Globals[g.Name]; exists {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Duplicate global value definition",
				Detail:   fmt.Sprintf("A global value named %q was already defined at %s. Global value names must be unique.", existing.Name, existing.DeclRange),
				Subject:  &g.DeclRange,
			})
		}
		m.Globals[g.Name] = g
	}

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

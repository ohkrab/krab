package configs

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/ohkrab/krab/diagnostics"
)

type Config struct {
	Module *Module
}

type Module struct {
	SourceDir     string
	Connections   map[string]*Connection
	Migrations    map[string]*Migration
	MigrationSets map[string]*MigrationSet
	// Globals     map[string]*Global
}

type File struct {
	Connections   []*Connection   `hcl:"connection"`
	Migrations    []*Migration    `hcl:"migration"`
	MigrationSets []*MigrationSet `hcl:"migration_set"`
	// Globals     []Global     `hcl:"globals"`
}

func NewModule(files []*File) (*Module, diagnostics.List) {
	diags := diagnostics.New()
	mod := &Module{
		Connections:   map[string]*Connection{},
		Migrations:    map[string]*Migration{},
		MigrationSets: map[string]*MigrationSet{},
		// Globals:     map[string]*Global{},
	}

	for _, file := range files {
		fileDiags := mod.appendFile(file)
		diags.Append(fileDiags)
	}

	return mod, diags
}

func (m *Module) appendFile(file *File) diagnostics.List {
	diags := diagnostics.New()

	// for _, g := range file.Globals {
	// 	if existing, exists := m.Globals[g.Name]; exists {
	// 		diags.Append(&hcl.Diagnostic{
	// 			Severity: hcl.DiagError,
	// 			Summary:  "Duplicate global value definition",
	// 			Detail:   fmt.Sprintf("A global value named %q was already defined at %s. Global value names must be unique.", existing.Name, existing.DeclRange),
	// 			Subject:  &g.DeclRange,
	// 		})
	// 	}
	// 	m.Globals[g.Name] = g
	// }

	for _, r := range file.Connections {
		key := r.Addr.String()
		if existing, exists := m.Connections[key]; exists {
			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("Duplicate resource %q configuration", existing.Type),
				Detail:   fmt.Sprintf("A %q resource named %q was already declared at %s. Resource names must be unique per type.", existing.Type, existing.Name, existing.DeclRange),
				Subject:  &r.DeclRange,
			})
			continue
		}
		m.Connections[key] = r
	}

	for _, r := range file.Migrations {
		key := r.Addr.String()
		if existing, exists := m.Migrations[key]; exists {
			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("Duplicate resource %q configuration", existing.Type),
				Detail:   fmt.Sprintf("A %q resource named %q was already declared at %s. Resource names must be unique per type.", existing.Type, existing.Name, existing.DeclRange),
				Subject:  &r.DeclRange,
			})
			continue
		}
		m.Migrations[key] = r
	}

	for _, r := range file.MigrationSets {
		key := r.Addr.String()
		if existing, exists := m.MigrationSets[key]; exists {
			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("Duplicate resource %q configuration", existing.Type),
				Detail:   fmt.Sprintf("A %q resource named %q was already declared at %s. Resource names must be unique per type.", existing.Type, existing.Name, existing.DeclRange),
				Subject:  &r.DeclRange,
			})
			continue
		}
		m.MigrationSets[key] = r
	}

	return diags
}

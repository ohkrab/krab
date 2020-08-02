package configs

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
)

func (p *Parser) LoadConfigDir(path string) (*Config, hcl.Diagnostics) {
	paths, diags := p.dirFiles(path)
	if diags.HasErrors() {
		return nil, diags
	}

	files, fDiags := p.loadFiles(paths)
	diags = append(diags, fDiags...)

	mod, modDiags := NewModule(files)
	diags = append(diags, modDiags...)

	mod.SourceDir = path

	config := &Config{
		Module: mod,
	}

	return config, diags
}

func (p *Parser) loadFiles(paths []string) ([]*File, hcl.Diagnostics) {
	var files []*File
	var diags hcl.Diagnostics

	for _, path := range paths {
		f, fDiags := p.LoadConfigFile(path)
		diags = append(diags, fDiags...)
		if f != nil {
			files = append(files, f)
		}
	}

	return files, diags
}

func (p *Parser) dirFiles(dir string) (primary []string, diags hcl.Diagnostics) {
	infos, err := p.fs.ReadDir(dir)
	if err != nil {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Failed to read module directory",
			Detail:   fmt.Sprintf("Module directory %s does not exist or cannot be read.", dir),
		})
		return
	}

	for _, info := range infos {
		if info.IsDir() {
			continue
		}

		name := info.Name()
		ext := fileExt(name)
		if ext == "" || isIgnoredFile(name) {
			continue
		}

		fullPath := filepath.Join(dir, name)
		primary = append(primary, fullPath)
	}

	return
}

func fileExt(path string) string {
	if strings.HasSuffix(path, ".krab") {
		return ".krab"
	} else {
		return "" // unrecognized extension
	}
}

func isIgnoredFile(name string) bool {
	return strings.HasPrefix(name, ".") || // Unix-like hidden files
		strings.HasSuffix(name, "~") || // vim
		strings.HasPrefix(name, "#") && strings.HasSuffix(name, "#") // emacs
}

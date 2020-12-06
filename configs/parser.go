package configs

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/ohkrab/krab/diagnostics"
	"github.com/spf13/afero"
)

// Parser is the main krab data parser.
type Parser struct {
	fs afero.Afero
	p  *hclparse.Parser
}

// SourceInfo contains information about the source for better debugging.
type SourceInfo struct {
	DeclRange hcl.Range
}

// NewParser creates krab parser.
func NewParser() *Parser {
	return &Parser{
		fs: afero.Afero{Fs: afero.OsFs{}},
		p:  hclparse.NewParser(),
	}
}

func (p *Parser) LoadConfigDir(path string) (*Config, diagnostics.List) {
	paths, diags := p.dirFiles(path)
	if diags.HasErrors() {
		return nil, diags
	}

	files, fDiags := p.loadConfigFiles(paths...)
	diags.Append(fDiags)

	mod, modDiags := NewModule(files)
	diags.Append(modDiags)

	mod.SourceDir = path

	config := &Config{
		Module: mod,
	}

	return config, diags
}

func (p *Parser) loadConfigFiles(paths ...string) ([]*File, diagnostics.List) {
	var files []*File
	diags := diagnostics.New()

	for _, path := range paths {
		f, fDiags := p.loadConfigFile(path)
		diags.Append(fDiags)
		if f != nil {
			files = append(files, f)
		}
	}

	return files, diags
}

func (p *Parser) loadConfigFile(path string) (*File, diagnostics.List) {
	var file File
	diags := diagnostics.New()

	if err := hclsimple.DecodeFile(path, nil, &file); err != nil {
		diags.Append(err)
		return nil, diags
	}

	return &file, diags
}

func (p *Parser) dirFiles(dir string) ([]string, diagnostics.List) {
	diags := diagnostics.New()
	paths := []string{}

	infos, err := p.fs.ReadDir(dir)
	if err != nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Failed to read directory",
			Detail:   fmt.Sprintf("Directory %s does not exist or cannot be read.", dir),
		})
		return paths, diags
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
		paths = append(paths, fullPath)
	}

	return paths, diags
}

func fileExt(path string) string {
	if strings.HasSuffix(path, ".krab.hcl") {
		return ".krab.hcl"
	} else {
		return "" // unrecognized extension
	}
}

func isIgnoredFile(name string) bool {
	return strings.HasPrefix(name, ".") || // Unix-like hidden files
		strings.HasSuffix(name, "~") || // vim
		strings.HasPrefix(name, "#") && strings.HasSuffix(name, "#") // emacs
}

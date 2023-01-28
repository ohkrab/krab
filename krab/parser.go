package krab

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl2/hclparse"
	"github.com/ohkrab/krab/krabenv"
	"github.com/ohkrab/krab/krabfn"
	"github.com/spf13/afero"
)

// Parser represents HCL simple parser.
type Parser struct {
	p  *hclparse.Parser
	FS afero.Afero
}

// NewParser initializes HCL parser and default file system.
func NewParser() *Parser {
	return &Parser{
		p:  hclparse.NewParser(),
		FS: afero.Afero{Fs: afero.OsFs{}},
	}
}

// LoadConfigDir parses files in a dir and returns Config.
func (p *Parser) LoadConfigDir(path string) (*Config, error) {
	paths, err := p.dirFiles(path)
	if err != nil {
		return nil, err
	}

	files, err := p.loadConfigFiles(paths...)
	if err != nil {
		return nil, err
	}

	cfg, err := NewConfig(files)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (p *Parser) loadConfigFiles(paths ...string) ([]*File, error) {
	var files []*File
	evalContext := krabfn.EvalContext(p.FS)

	for _, path := range paths {
		f, err := p.loadConfigFile(path, evalContext)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}

	return files, nil
}

func (p *Parser) loadConfigFile(path string, evalContext *hcl.EvalContext) (*File, error) {
	src, err := p.FS.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("[%w] Failed to load file %s", err, path)
	}

	hclFile, diags := hclsyntax.ParseConfig(src, path, hcl.Pos{Line: 1, Column: 1, Byte: 0})
	if diags.HasErrors() {
		return nil, fmt.Errorf("[%s] Failed to decode file %s", err.Error(), path)
	}
	file := &File{File: hclFile}
	if err := file.Decode(evalContext); err != nil {
		return nil, fmt.Errorf("[%w] Failed to decode file %s", err, path)
	}

	return file, nil
}

func (p *Parser) dirFiles(dir string) ([]string, error) {
	paths := []string{}

	infos, err := p.FS.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("Directory %s does not exist or cannot be read", dir)
	}

	for _, info := range infos {
		if info.IsDir() {
			fullPath := filepath.Join(dir, info.Name())
			nestedPaths, err := p.dirFiles(fullPath)
			if err != nil {
				return nil, err
			}
			paths = append(paths, nestedPaths...)
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

	return paths, nil
}

func fileExt(path string) string {
	if strings.HasSuffix(path, krabenv.Ext()) {
		return krabenv.Ext()
	}

	return "" // unrecognized
}

func isIgnoredFile(name string) bool {
	return strings.HasPrefix(name, ".") || // dotfiles
		strings.HasSuffix(name, "~") || // vim/backups
		strings.HasPrefix(name, "#") && strings.HasSuffix(name, "#") // emacs
}

package krab

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/hcl2/hclparse"
	"github.com/spf13/afero"
)

type Parser struct {
	p  *hclparse.Parser
	fs afero.Afero
}

// NewParser initializes HCL parser and default file system.
func NewParser() *Parser {
	return &Parser{
		p:  hclparse.NewParser(),
		fs: afero.Afero{Fs: afero.OsFs{}},
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

	for _, path := range paths {
		f, err := p.loadConfigFile(path)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}

	return files, nil
}

func (p *Parser) loadConfigFile(path string) (*File, error) {
	var file File

	src, err := p.fs.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("[%w] Failed to load file %s", err, path)
	}

	if err := hclsimple.Decode(path, src, nil, &file); err != nil {
		return nil, fmt.Errorf("[%w] Failed to decode file %s", err, path)
	}

	return &file, nil
}

func (p *Parser) dirFiles(dir string) ([]string, error) {
	paths := []string{}

	infos, err := p.fs.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("Directory %s does not exist or cannot be read.", dir)
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

	return paths, nil
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

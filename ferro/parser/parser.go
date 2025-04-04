package parser

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/ohkrab/krab/ferro/config"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
)

const (
	FileExt = ".fyml"
)

type Parser struct {
	FS fs.FS
}

// New initializes parser and default file system.
func New() *Parser {
	return &Parser{
		FS: os.DirFS("."),
	}
}

// LoadConfigDir parses files in a dir and returns ParsedConfig.
func (p *Parser) LoadConfigDir(path string) (*ParsedConfig, error) {
	paths, err := p.dirFiles(path)
	if err != nil {
		return nil, err
	}

	parsedFiles, err := p.loadParsedFiles(paths)
	if err != nil {
		return nil, err
	}

	if err := p.parse(parsedFiles); err != nil {
		return nil, err
	}

	cfg := &ParsedConfig{
		Files: parsedFiles,
	}

	return cfg, nil
}

func (p *Parser) parse(files []*ParsedFile) error {
	for _, file := range files {
		for _, chunk := range file.Chunks {
			var parseErr error
			switch chunk.Header.ApiVersion {
			case "migrations/v1":
				parseErr = p.parseMigrationsV1(file, chunk)
			default:
				parseErr = fmt.Errorf("unsupported api version: %s (%s)", chunk.Header.ApiVersion, file.Path)
			}

			if parseErr != nil {
				return parseErr
			}
		}
	}

	return nil
}

func (p *Parser) parseMigrationsV1(file *ParsedFile, chunk *ParsedChunk) error {
	switch chunk.Header.Kind {
	case "Migration":
		var migration config.Migration
		if err := yaml.Unmarshal(chunk.Raw, &migration); err != nil {
			return fmt.Errorf("failed to parse Migration: %w", err)
		}
		file.Migrations = append(file.Migrations, &migration)

	case "MigrationSet":
		var migrationSet config.MigrationSet
		if err := yaml.Unmarshal(chunk.Raw, &migrationSet); err != nil {
			return fmt.Errorf("failed to parse MigrationSet: %w", err)
		}
		file.MigrationSets = append(file.MigrationSets, &migrationSet)

	default:
		return fmt.Errorf("unsupported kind: %s (%s)", chunk.Header.Kind, file.Path)
	}

	return nil
}

func (p *Parser) loadParsedFiles(paths []string) ([]*ParsedFile, error) {
	parsedFiles := make([]*ParsedFile, len(paths))
	for i, path := range paths {
		parsedFiles[i] = &ParsedFile{
			Path:   path,
			Chunks: []*ParsedChunk{},
		}
	}

	eg := errgroup.Group{}
	for f, file := range parsedFiles {
		eg.Go(func() error {
			src, err := fs.ReadFile(p.FS, file.Path)
			if err != nil {
				return fmt.Errorf("failed to read config file %s: %w", file.Path, err)
			}

			chunks := bytes.Split(src, []byte("\n---"))
			for i, chunk := range chunks {
				var parsedHeader config.Header
				if err := yaml.Unmarshal(chunk, &parsedHeader); err != nil {
					return fmt.Errorf("failed to unmarshal chunk (%d): %w", i, err)
				}

				parsedFiles[f].Chunks = append(parsedFiles[f].Chunks, &ParsedChunk{
					Header: &parsedHeader,
					Raw:    chunk,
				})
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return parsedFiles, nil
}

func (p *Parser) dirFiles(root string) ([]string, error) {
	paths := []string{}

	err := fs.WalkDir(p.FS, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		ext := fileExt(path)
		if ext == "" || isIgnoredFile(path) {
			return nil
		}

		paths = append(paths, path)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return paths, nil
}

func fileExt(path string) string {
	if strings.HasSuffix(path, FileExt) {
		return FileExt
	}

	return "" // unrecognized
}

func isIgnoredFile(name string) bool {
	return strings.HasPrefix(name, ".") || // dotfiles
		strings.HasSuffix(name, "~") || // vim/backups
		strings.HasPrefix(name, "#") && strings.HasSuffix(name, "#") // emacs
}

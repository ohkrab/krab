package parser

import (
	"fmt"
	"path/filepath"

	"github.com/ohkrab/krab/ferro/config"
	"gopkg.in/yaml.v3"
)

type Parser struct {
	fs *config.Filesystem
}

// New initializes parser and default file system.
func New(fs *config.Filesystem) *Parser {
	return &Parser{
		fs: fs,
	}
}

// LoadAndParse parses files in a dir and returns ParsedConfig.
func (p *Parser) LoadAndParse() (*config.ParsedConfig, error) {
	paths, err := p.fs.DirFiles()
	if err != nil {
		return nil, err
	}

	parsedFiles, err := p.fs.LoadFiles(paths)
	if err != nil {
		return nil, err
	}

	if err := p.parse(parsedFiles); err != nil {
		return nil, err
	}

	cfg := &config.ParsedConfig{
		Files: parsedFiles,
	}

	return cfg, nil
}

func (p *Parser) parse(files []*config.ParsedFile) error {
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

func (p *Parser) parseMigrationsV1(file *config.ParsedFile, chunk *config.ParsedChunk) error {
	switch chunk.Header.Kind {
	case "Migration":
		var migration config.Migration
		if err := yaml.Unmarshal(chunk.Raw, &migration); err != nil {
			return fmt.Errorf("failed to parse Migration: %w", err)
		}
		migration.Path = filepath.Join(p.fs.Dir, file.Path)
		file.Migrations = append(file.Migrations, &migration)

	case "MigrationSet":
		var migrationSet config.MigrationSet
		if err := yaml.Unmarshal(chunk.Raw, &migrationSet); err != nil {
			return fmt.Errorf("failed to parse MigrationSet: %w", err)
		}
		migrationSet.Path = filepath.Join(p.fs.Dir, file.Path)
		file.MigrationSets = append(file.MigrationSets, &migrationSet)

	default:
		return fmt.Errorf("unsupported kind: %s (%s)", chunk.Header.Kind, file.Path)
	}

	return nil
}

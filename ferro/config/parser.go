package config

import (
	"fmt"
	"path/filepath"

	"github.com/ohkrab/krab/fmtx"
	"gopkg.in/yaml.v3"
)

type Parser struct {
	fs     *Filesystem
	logger *fmtx.Logger
}

// New initializes parser and default file system.
func NewParser(fs *Filesystem, logger *fmtx.Logger) *Parser {
	return &Parser{
		fs: fs,
        logger: logger,
	}
}

// LoadAndParse parses files in a dir and returns ParsedConfig.
func (p *Parser) LoadAndParse() (*ParsedConfig, error) {
	paths, err := p.fs.DirFiles()
	if err != nil {
		return nil, err
	}

	parsedFiles, err := p.fs.LoadFiles(paths)
	if err != nil {
		return nil, err
	}

	for _, file := range parsedFiles {
		p.logger.WriteInfo("  using file: %s", file.Path)
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
			case "drivers/v1":
				parseErr = p.parseDriversV1(file, chunk)
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
		var migration Migration
		if err := yaml.Unmarshal(chunk.Raw, &migration); err != nil {
			return fmt.Errorf("failed to parse Migration: %w\n%s", err, string(chunk.Raw))
		}
		migration.Path = filepath.Join(p.fs.Dir, file.Path)
		file.Migrations = append(file.Migrations, &migration)

	case "MigrationSet":
		var migrationSet MigrationSet
		if err := yaml.Unmarshal(chunk.Raw, &migrationSet); err != nil {
			return fmt.Errorf("failed to parse MigrationSet: %w\n%s", err, string(chunk.Raw))
		}
		migrationSet.Path = filepath.Join(p.fs.Dir, file.Path)
		file.MigrationSets = append(file.MigrationSets, &migrationSet)

	default:
		return fmt.Errorf("unsupported kind: %s (%s)", chunk.Header.Kind, file.Path)
	}

	return nil
}

func (p *Parser) parseDriversV1(file *ParsedFile, chunk *ParsedChunk) error {
	switch chunk.Header.Kind {
	case "Driver":
		var driver Driver
		if err := yaml.Unmarshal(chunk.Raw, &driver); err != nil {
			return fmt.Errorf("failed to parse Driver: %w\n%s", err, string(chunk.Raw))
		}
		driver.Path = filepath.Join(p.fs.Dir, file.Path)
		file.Drivers = append(file.Drivers, &driver)

	default:
		return fmt.Errorf("unsupported kind: %s (%s)", chunk.Header.Kind, file.Path)
	}

	return nil
}

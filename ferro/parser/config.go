package parser

import (
	"fmt"

	"github.com/ohkrab/krab/ferro/config"
)

type ParsedConfig struct {
	Files []*ParsedFile
}

type ParsedFile struct {
	Path   string
	Chunks []*ParsedChunk

	Migrations    []*config.Migration
	MigrationSets []*config.MigrationSet
}

type ParsedChunk struct {
	Header *config.Header
	Raw    []byte
}

func (p *ParsedConfig) BuildConfig() (*config.Config, error) {
	cfg := config.New()

	for _, file := range p.Files {
		for _, migration := range file.Migrations {
			if err := cfg.AddMigration(migration); err != nil {
				return nil, fmt.Errorf("adding Migration: %w", err)
			}
		}
		for _, migrationSet := range file.MigrationSets {
			if err := cfg.AddMigrationSet(migrationSet); err != nil {
				return nil, fmt.Errorf("adding MigrationSet: %w", err)
			}
		}
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}

	return cfg, nil
}

package run

import (
	"fmt"

	"github.com/ohkrab/krab/ferro/config"
	"github.com/ohkrab/krab/plugins"
)

type Builder struct {
	fs       *config.Filesystem
	parsed   *config.ParsedConfig
	registry *plugins.Registry
}

func NewBuilder(
	fs *config.Filesystem,
	parsed *config.ParsedConfig,
	registry *plugins.Registry,
) *Builder {
	return &Builder{fs: fs, parsed: parsed, registry: registry}
}

func (b *Builder) Validate(cfg *config.Config) *config.Errors {
	errors := &config.Errors{
		Errors: []error{},
	}

	//// validate resources
	// Migrations
	for _, resource := range cfg.Migrations {
		if errs := resource.Validate(); errs != nil {
			for _, err := range errs.Errors {
				errors.Append(err)
			}
		}
	}

	// MigrationSets
	for _, resource := range cfg.MigrationSets {
		if errs := resource.Validate(); errs != nil {
			for _, err := range errs.Errors {
				errors.Append(err)
			}
		}
	}

	// Drivers
	for _, resource := range cfg.Drivers {
		if errs := resource.Validate(); errs != nil {
			for _, err := range errs.Errors {
				errors.Append(err)
			}
		}
	}

	//// post-validation when all resources are loaded
	// check invalid references
	for _, migrationSet := range cfg.MigrationSets {
		for _, migrationName := range migrationSet.Spec.Migrations {
			if _, ok := cfg.Migrations[migrationName]; !ok {
				errors.Append(
					fmt.Errorf("invalid reference: Migration `%s` (referenced by MigrationSet `%s`) does not exist", migrationName, migrationSet.Metadata.Name),
				)
			}
		}
	}

	return errors
}

func (b *Builder) BuildConfig() (*config.Config, *config.Errors) {
	cfg := config.New()

	for _, file := range b.parsed.Files {
		for _, migration := range file.Migrations {
			if err := cfg.AddMigration(migration); err != nil {
				return nil, config.Errorf("adding Migration: %w", err)
			}
		}
		for _, migrationSet := range file.MigrationSets {
			if err := cfg.AddMigrationSet(migrationSet); err != nil {
				return nil, config.Errorf("adding MigrationSet: %w", err)
			}
		}
		for _, driver := range file.Drivers {
			if err := cfg.AddDriver(driver); err != nil {
				return nil, config.Errorf("adding Driver: %w", err)
			}
		}
	}

	errors := b.Validate(cfg)
	return cfg, errors
}

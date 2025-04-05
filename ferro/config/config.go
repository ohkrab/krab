package config

import "fmt"

type Config struct {
	MigrationSets map[string]*MigrationSet
	Migrations    map[string]*Migration
}

type Resource interface {
	EnforceDefaults()
	Validate() *Errors
}

func New() *Config {
	return &Config{
		MigrationSets: make(map[string]*MigrationSet),
		Migrations:    make(map[string]*Migration),
	}
}

func (c *Config) AddMigrationSet(migrationSet *MigrationSet) error {
	if _, ok := c.MigrationSets[migrationSet.Metadata.Name]; ok {
		return fmt.Errorf("migration set `%s` already exists", migrationSet.Metadata.Name)
	}
	c.MigrationSets[migrationSet.Metadata.Name] = migrationSet
	migrationSet.EnforceDefaults()

	return nil
}

func (c *Config) AddMigration(migration *Migration) error {
	if _, ok := c.Migrations[migration.Metadata.Name]; ok {
		return fmt.Errorf("migration `%s` already exists", migration.Metadata.Name)
	}
	c.Migrations[migration.Metadata.Name] = migration
	migration.EnforceDefaults()

	return nil
}

func (c *Config) Validate() *Errors {
	errors := &Errors{
		Errors: []error{},
	}
	// check invalid references
	for _, migrationSet := range c.MigrationSets {
		for _, migrationName := range migrationSet.Spec.Migrations {
			if _, ok := c.Migrations[migrationName]; !ok {
				errors.Append(
					fmt.Errorf("invalid reference: Migration `%s` (referenced by MigrationSet `%s`) does not exist", migrationName, migrationSet.Metadata.Name),
				)
			}
		}
	}

	// validate resources
	for _, resource := range c.Migrations {
		if errs := resource.Validate(); errs != nil {
			for _, err := range errs.Errors {
				errors.Append(err)
			}
		}
	}

	if len(errors.Errors) > 0 {
		return errors
	}

	return nil
}

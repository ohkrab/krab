package config

import "fmt"

type Config struct {
	MigrationSets map[string]*MigrationSet
	Migrations    map[string]*Migration
}

type Resource interface {
	EnforceDefaults()
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

func (c *Config) Validate() error {
	return nil
}

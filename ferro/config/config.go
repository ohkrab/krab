package config

import (
	"fmt"
)

type Config struct {
	MigrationSets map[string]*MigrationSet
	Migrations    map[string]*Migration
	Drivers       map[string]*Driver
}

type Resource interface {
	EnforceDefaults()
	Validate() *Errors
}

func New() *Config {
	return &Config{
		MigrationSets: make(map[string]*MigrationSet),
		Migrations:    make(map[string]*Migration),
		Drivers:       make(map[string]*Driver),
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

	err := migration.ResolveFiles()
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) AddDriver(driver *Driver) error {
	if _, ok := c.Drivers[driver.Metadata.Name]; ok {
		return fmt.Errorf("driver `%s` already exists", driver.Metadata.Name)
	}
	c.Drivers[driver.Metadata.Name] = driver
	driver.EnforceDefaults()

	return nil
}

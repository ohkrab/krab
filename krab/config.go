package krab

import (
	"fmt"

	"github.com/ohkrab/krab/krabhcl"
)

// Config represents all configuration loaded from directory.
//
type Config struct {
	MigrationSets map[string]*MigrationSet
	Migrations    map[string]*Migration
}

// NewConfig returns new configuration that was read from Parser.
// Transient attributes are updated with parsed data.
func NewConfig(files []*File) (*Config, error) {
	c := &Config{
		MigrationSets: make(map[string]*MigrationSet),
		Migrations:    make(map[string]*Migration),
	}

	// append files
	for _, f := range files {
		if err := c.appendFile(f); err != nil {
			return nil, err
		}
	}

	// parse refs
	for _, set := range c.MigrationSets {
		set.Migrations = make([]*Migration, 0)

		traversals := set.MigrationsExpr.Variables()
		for _, t := range traversals {
			addr, err := parseTraversalToAddr(t)
			if err != nil {
				return nil, fmt.Errorf("Parsing migrations for set '%s' failed. %w", set.RefName, err)
			}
			migration, found := c.Migrations[addr.OnlyRefNames()]
			if !found {
				return nil, fmt.Errorf("Migration Set references '%s' migration that does not exist", addr.OnlyRefNames())
			}
			set.Migrations = append(set.Migrations, migration)
		}
	}

	// defaults
	for _, defaultable := range c.MigrationSets {
		defaultable.InitDefaults()
	}

	// validate
	for _, validatable := range c.MigrationSets {
		if err := validatable.Validate(); err != nil {
			return nil, err
		}
	}

	for _, validatable := range c.Migrations {
		if err := validatable.Validate(); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *Config) appendFile(file *File) error {
	for i, m := range file.Migrations {
		if _, found := c.Migrations[m.RefName]; found {
			return fmt.Errorf("Migration with the name '%s' already exists", m.RefName)
		}

		raw := file.Raw.Migrations[i]
		c.Migrations[m.RefName] = m

		applyDefRangesToMigration(&m.Up, &raw.Up)
		applyDefRangesToMigration(&m.Down, &raw.Down)
	}

	for _, s := range file.MigrationSets {
		if _, found := c.MigrationSets[s.RefName]; found {
			return fmt.Errorf("Migration Set with the name '%s' already exists", s.RefName)
		}

		c.MigrationSets[s.RefName] = s
	}

	return nil
}

func applyDefRangesToMigration(m *MigrationUpOrDown, raw *RawMigrationUpOrDown) {
	remain := krabhcl.Body{raw.Remain}
	for i, defRange := range remain.DefRangesFromPartialContent(&DDLCreateTableSchema) {
		m.CreateTables[i].DefRange = defRange
	}
	for i, defRange := range remain.DefRangesFromPartialContent(&DDLDropTableSchema) {
		m.DropTables[i].DefRange = defRange
	}
	for i, defRange := range remain.DefRangesFromPartialContent(&DDLCreateIndexSchema) {
		m.CreateIndices[i].DefRange = defRange
	}
	for i, defRange := range remain.DefRangesFromPartialContent(&DDLDropIndexSchema) {
		m.DropIndices[i].DefRange = defRange
	}
}

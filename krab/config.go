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
	Actions       map[string]*Action
	Wasms         map[string]*WebAssembly
}

// NewConfig returns new configuration that was read from Parser.
// Transient attributes are updated with parsed data.
func NewConfig(files []*File) (*Config, error) {
	c := &Config{
		MigrationSets: map[string]*MigrationSet{},
		Migrations:    map[string]*Migration{},
		Actions:       map[string]*Action{},
		Wasms:         map[string]*WebAssembly{},
	}

	// append files
	for _, f := range files {
		if err := c.appendFile(f); err != nil {
			return nil, err
		}
	}

	// parse refs
	for _, set := range c.MigrationSets {
		set.Migrations = []*Migration{}

		traversals := set.MigrationsExpr.Variables()
		for _, t := range traversals {
			addr, err := krabhcl.ParseTraversalToAddr(t)
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
	for _, defaultable := range c.Actions {
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

	for _, a := range file.Actions {
		if _, found := c.Actions[a.Addr().OnlyRefNames()]; found {
			return fmt.Errorf("Action with the name '%s' already exists", a.Addr().OnlyRefNames())
		}

		c.Actions[a.Addr().OnlyRefNames()] = a
	}

	return nil
}

func applyDefRangesToMigration(m *MigrationUpOrDown, raw *RawMigrationUpOrDown) {
	remain := krabhcl.Body{raw.Remain}

	for i, defRange := range remain.DefRangesFromPartialContentBlocks(&DDLCreateTableSchema) {
		m.CreateTables[i].DefRange = defRange
	}
	for i, defRange := range remain.DefRangesFromPartialContentBlocks(&DDLDropTableSchema) {
		m.DropTables[i].DefRange = defRange
	}
	for i, defRange := range remain.DefRangesFromPartialContentBlocks(&DDLCreateIndexSchema) {
		m.CreateIndices[i].DefRange = defRange
	}
	for i, defRange := range remain.DefRangesFromPartialContentBlocks(&DDLDropIndexSchema) {
		m.DropIndices[i].DefRange = defRange
	}

	m.AttrDefRanges = remain.DefRangesFromPartialContentAttributes(&MigrationUpOrDownSchema)
}

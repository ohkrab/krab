package krab

import (
	"fmt"
)

// Config represents all configuration loaded from directory.
type Config struct {
	MigrationSets map[string]*MigrationSet
	Migrations    map[string]*Migration
	Actions       map[string]*Action
	TestSuite     *TestSuite
	TestExamples  map[string]*TestExample
}

// NewConfig returns new configuration that was read from Parser.
// Transient attributes are updated with parsed data.
func NewConfig(files []*File) (*Config, error) {
	c := &Config{
		MigrationSets: map[string]*MigrationSet{},
		Migrations:    map[string]*Migration{},
		Actions:       map[string]*Action{},
		TestSuite:     &TestSuite{},
		TestExamples:  map[string]*TestExample{},
	}

	// append files
	for _, f := range files {
		if err := c.appendFile(f); err != nil {
			return nil, err
		}
	}

	// connect by refs
	for _, set := range c.MigrationSets {
		for _, addr := range set.MigrationAddrs {
			migration, found := c.Migrations[addr.OnlyRefNames()]
			if !found {
				return nil, fmt.Errorf("Migration Set references '%s' migration that does not exist", addr.OnlyRefNames())
			}
			set.Migrations = append(set.Migrations, migration)
		}
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

	for _, validatable := range c.Actions {
		if err := validatable.Validate(); err != nil {
			return nil, err
		}
	}

	for _, validatable := range c.TestExamples {
		if err := validatable.Validate(); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *Config) appendFile(file *File) error {
	for _, m := range file.Migrations {
		if _, found := c.Migrations[m.RefName]; found {
			return fmt.Errorf("Migration with the name '%s' already exists", m.RefName)
		}

		c.Migrations[m.RefName] = m
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

	c.TestSuite = file.TestSuite
	c.TestSuite.Tests = []*TestExample{}

	for _, t := range file.TestExamples {
		if _, found := c.TestExamples[t.Addr().OnlyRefNames()]; found {
			return fmt.Errorf("Test with the name '%s' already exists", t.Addr().OnlyRefNames())
		}
		c.TestSuite.Tests = append(c.TestSuite.Tests, t)
	}

	return nil
}

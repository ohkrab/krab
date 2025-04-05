package config

import "fmt"

// Migration

type Migration struct {
	Metadata Metadata      `yaml:"metadata"`
	Spec     MigrationSpec `yaml:"spec"`
}

type MigrationSpec struct {
	Version     string       `yaml:"version"`
	Transaction *bool        `yaml:"transaction,omitempty"`
	Run         MigrationRun `yaml:"run"`
}

type MigrationRun struct {
	Up struct {
		Sql  string `yaml:"sql"`
		File string `yaml:"file"`
	} `yaml:"up"`
	Down struct {
		Sql  string `yaml:"sql"`
		File string `yaml:"file"`
	} `yaml:"down"`
}

func (m *Migration) EnforceDefaults() {
	if m.Spec.Transaction == nil {
		defaultTransaction := true
		m.Spec.Transaction = &defaultTransaction
	}
}

func (m *Migration) Validate() *Errors {
	errors := &Errors{
		Errors: []error{},
	}

	// break early if no name because it's hard to explain other errors without it
	if m.Metadata.Name == "" {
		errors.Append(fmt.Errorf("invalid spec: Migration must have a name"))
		return errors
	}

	// sanity check
	if m.Spec.Transaction == nil {
		panic(fmt.Sprintf("invalid spec: Migration(transaction) `%s` transaction is nil", m.Metadata.Name))
	}

	if m.Spec.Run.Up.Sql != "" && m.Spec.Run.Up.File != "" {
		errors.Append(fmt.Errorf("invalid spec: Migration(up) `%s` cannot have both `sql` and `file` defined", m.Metadata.Name))
	}

	if m.Spec.Run.Down.Sql != "" && m.Spec.Run.Down.File != "" {
		errors.Append(fmt.Errorf("invalid spec: Migration(down) `%s` cannot have both `sql` and `file` defined", m.Metadata.Name))
	}

	if m.Spec.Run.Up.Sql == "" && m.Spec.Run.Up.File == "" {
		errors.Append(fmt.Errorf("invalid spec: Migration(up) `%s` must have either `sql` or `file` defined", m.Metadata.Name))
	}

	if m.Spec.Run.Down.Sql == "" && m.Spec.Run.Down.File == "" {
		errors.Append(fmt.Errorf("invalid spec: Migration(down) `%s` must have either `sql` or `file` defined", m.Metadata.Name))
	}

	return errors
}

// MigrationSet

type MigrationSet struct {
	Metadata Metadata         `yaml:"metadata"`
	Spec     MigrationSetSpec `yaml:"spec"`
}

type MigrationSetSpec struct {
	Migrations []string `yaml:"migrations"`
}

func (m *MigrationSet) EnforceDefaults() {
	if m.Spec.Migrations == nil {
		m.Spec.Migrations = []string{}
	}
}

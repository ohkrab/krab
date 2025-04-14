package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Migration

type Migration struct {
	Path     string        `yaml:"-"`
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

func (m *Migration) ResolveFiles() error {
	if m.Spec.Run.Up.File != "" {
		path := m.Spec.Run.Up.File
		if !filepath.IsAbs(path) {
			path = filepath.Join(filepath.Dir(m.Path), path)
		}

		upFile, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("io error: Migration(up) `%s` cannot load file `%s`", m.Metadata.Name, path)
		}
		defer upFile.Close()

		upBytes, err := io.ReadAll(upFile)
		if err != nil {
			return fmt.Errorf("io error: Migration(up) `%s` cannot read file `%s`", m.Metadata.Name, path)
		}
		m.Spec.Run.Up.Sql = string(upBytes)
	}

	if m.Spec.Run.Down.File != "" {
		path := m.Spec.Run.Down.File
		if !filepath.IsAbs(path) {
			path = filepath.Join(filepath.Dir(m.Path), path)
		}

		downFile, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("io error: Migration(down) `%s` cannot load file `%s`", m.Metadata.Name, m.Spec.Run.Down.File)
		}
		defer downFile.Close()

		downBytes, err := io.ReadAll(downFile)
		if err != nil {
			return fmt.Errorf("io error: Migration(down) `%s` cannot read file `%s`", m.Metadata.Name, m.Spec.Run.Down.File)
		}
		m.Spec.Run.Down.Sql = string(downBytes)
	}

	return nil
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

	if !errors.HasErrors() {
		err := m.ResolveFiles()
		if err != nil {
			errors.Append(err)
		}
	}

	return errors
}

// MigrationSet

type MigrationSet struct {
	Path     string           `yaml:"-"`
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

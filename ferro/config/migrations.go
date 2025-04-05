package config

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
		Sql string `yaml:"sql"`
	} `yaml:"up"`
	Down struct {
		Sql string `yaml:"sql"`
	} `yaml:"down"`
}

func (m *Migration) EnforceDefaults() {
	if m.Spec.Transaction == nil {
		defaultTransaction := true
		m.Spec.Transaction = &defaultTransaction
	}
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

package krab

// Migration represents single up/down migration pair.
//
type Migration struct {
	RefName string `hcl:"ref_name,label"`

	Version string        `hcl:"version,optional"`
	Up      MigrationUp   `hcl:"up,block"`
	Down    MigrationDown `hcl:"down,block"`
}

// OnAfterParse updates version if not specified.
func (m *Migration) OnAfterParse() {
	if m.Version == "" {
		m.Version = m.RefName
	}
}

// MigrationUp contains info how to migrate up.
type MigrationUp struct {
	SQL string `hcl:"sql,optional"`
}

// MigrationDown contains info how to migrate down.
type MigrationDown struct {
	SQL string `hcl:"sql,optional"`
}

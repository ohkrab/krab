package krab

// Migration represents single up/down migration pair.
//
type Migration struct {
	RefName string `hcl:"ref_name,label"`

	Up   MigrationUp   `hcl:"up,block"`
	Down MigrationDown `hcl:"down,block"`
}

// MigrationUp contains info how to migrate up.
type MigrationUp struct {
	Sql string `hcl:"sql,optional"`
}

// MigrationUp contains info how to migrate down.
type MigrationDown struct {
	Sql string `hcl:"sql,optional"`
}

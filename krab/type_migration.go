package krab

// Migration represents single up/down migration pair.
//
type Migration struct {
	RefName string `hcl:"name,label"`

	Up   MigrationUp   `hcl:"up,block"`
	Down MigrationDown `hcl:"down,block"`
}

// MigrationUp contains info how to migrate up.
type MigrationUp struct {
	Sql string `hcl:"sql"`
}

// MigrationUp contains info how to migrate down.
type MigrationDown struct {
	Sql string `hcl:"sql"`
}

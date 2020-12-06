package configs

import (
	"github.com/ohkrab/krab/addrs"
)

// Migration represents single up/down migration pair.
type Migration struct {
	addrs.Addr
	SourceInfo

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

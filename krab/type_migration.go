package krab

import "fmt"

// Migration represents single up/down migration pair.
//
type Migration struct {
	RefName string `hcl:"ref_name,label"`

	Version     string        `hcl:"version"`
	Up          MigrationUp   `hcl:"up,block"`
	Down        MigrationDown `hcl:"down,block"`
	Transaction *bool         `hcl:"transaction,optional"` // wrap operation in transaction
}

// MigrationUp contains info how to migrate up.
type MigrationUp struct {
	SQL string `hcl:"sql,optional"`
}

// MigrationDown contains info how to migrate down.
type MigrationDown struct {
	SQL string `hcl:"sql,optional"`
}

func (ms *Migration) Validate() error {
	return ErrorCoalesce(
		ValidateRefName(ms.RefName),
		ValidateStringNonEmpty(fmt.Sprint("`version` attribute in `", ms.RefName, "` migration"), ms.Version),
	)
}

// ShouldRunInTransaction returns whether migration should be wrapped into transaction or not.
func (ms *Migration) ShouldRunInTransaction() bool {
	if ms.Transaction == nil {
		return true
	}
	return *ms.Transaction
}

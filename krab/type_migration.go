package krab

import (
	"fmt"
)

// Migration represents single up/down migration pair.
//
type Migration struct {
	RefName string `hcl:"ref_name,label"`

	Version     string            `hcl:"version"`
	Up          MigrationUpOrDown `hcl:"up,block"`
	Down        MigrationUpOrDown `hcl:"down,block"`
	Transaction *bool             `hcl:"transaction,optional"` // wrap operaiton in transaction
}

// Migration contains info how to migrate up or down.
type MigrationUpOrDown struct {
	SQL          string            `hcl:"sql,optional"`
	CreateTables []*DDLCreateTable `hcl:"create_table,block"`
	DropTables   []*DDLDropTable   `hcl:"drop_table,block"`
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

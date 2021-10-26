package krab

import "fmt"

// Migration represents single up/down migration pair.
//
type Migration struct {
	RefName string `hcl:"ref_name,label"`

	Version     string        `hcl:"version"`
	Up          MigrationUp   `hcl:"up,block"`
	Down        MigrationDown `hcl:"down,block"`
	Transaction *bool         `hcl:"transaction,optional"` // wrap operaiton in transaction
}

// MigrationUp contains info how to migrate up.
type MigrationUp struct {
	SQL string `hcl:"sql,optional"`
}

// MigrationDown contains info how to migrate down.
type MigrationDown struct {
	SQL string `hcl:"sql,optional"`
}

// CreateTable contains DSL for creating tables.
type CreateTable struct {
	Name     string `hcl:"name,label"`
	Unlogged bool
}

// PrimaryKey constraint DSL for table DDL.
type PrimaryKey struct {
	Name string `hcl:"name,label"`
}

// ForeignKey constraint DSL for table DDL.
type ForeignKey struct {
	Name string `hcl:"name,label"`
}

// ForeignKey constraint DSL for table DDL.
type Unique struct {
	Name string `hcl:"name,label"`
}

// Column constraint DSL for table DDL.
type Column struct {
	Name string `hcl:"name,label"`
	Type string `hcl:"type,label"`
}

// Check constraint DSL for table DDL.
type Check struct {
	Name string `hcl:"name,label"`
}

// DropTable contains DSL for dropping tables.
type DropTable struct {
	Table string `hcl:"table,label"`
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

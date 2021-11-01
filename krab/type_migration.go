package krab

import (
	"fmt"
	"io"

	"github.com/hashicorp/hcl/v2"
)

// ToSQL converts DSL struct to SQL.
type ToSQL interface {
	ToSQL(w io.StringWriter)
}

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
		ms.Up.Validate(),
		ms.Down.Validate(),
	)
}

// ShouldRunInTransaction returns whether migration should be wrapped into transaction or not.
func (ms *Migration) ShouldRunInTransaction() bool {
	if ms.Transaction == nil {
		return true
	}
	return *ms.Transaction
}

func (m *MigrationUpOrDown) Validate() error {
	return nil
}

// ToSQL converts migration definition to SQL.
func (m *MigrationUpOrDown) ToSQL(w io.StringWriter) {
	w.WriteString(m.SQL)

	for _, t := range m.CreateTables {
		t.ToSQL(w)
	}
	for _, t := range m.DropTables {
		t.ToSQL(w)
	}
}

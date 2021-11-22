package krab

import (
	"fmt"
	"io"
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
	SQL           string            `hcl:"sql,optional"`
	CreateTables  []*DDLCreateTable `hcl:"create_table,block"`
	CreateIndices []*DDLCreateIndex `hcl:"create_index,block"`
	DropTables    []*DDLDropTable   `hcl:"drop_table,block"`
	DropIndices   []*DDLDropIndex   `hcl:"drop_index,block"`
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

func (m *MigrationUpOrDown) ToSQL(w io.StringWriter) {
	w.WriteString(m.SQL)
}

// ToSQLStatements returns list of SQL statements to executre during the migration.
func (m *MigrationUpOrDown) ToSQLStatements() SQLStatements {
	sqls := SQLStatements{}

	if m.SQL != "" {
		sqls.Append(m)
	}

	for _, t := range m.CreateTables {
		sqls.Append(t)
	}
	for _, t := range m.CreateIndices {
		sqls.Append(t)
	}
	for _, t := range m.DropIndices {
		sqls.Append(t)
	}
	for _, t := range m.DropTables {
		sqls.Append(t)
	}

	return sqls
}

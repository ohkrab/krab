package krab

import (
	"github.com/hashicorp/hcl/v2"
)

// MigrationSet represents collection of migrations.
type MigrationSet struct {
	RefName string `hcl:"ref_name,label"`

	// SchemaMigrationsTable string         `hcl:"schema_migrations_table"`
	MigrationsExpr hcl.Expression `hcl:"migrations"`
	Migrations     []*Migration   // populated from refs in expression
}

func (ms *MigrationSet) Validate() error {
	return ErrorCoalesce(
		ValidateRefName(ms.RefName),
	)
}

// FindMigrationByVersion looks up for the migration in current set.
func (ms *MigrationSet) FindMigrationByVersion(version string) *Migration {
	for _, m := range ms.Migrations {
		if m.Version == version {
			return m
		}
	}

	return nil
}

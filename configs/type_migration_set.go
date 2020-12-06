package configs

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/ohkrab/krab/addrs"
)

// MigrationSet represents collection of migrations.
type MigrationSet struct {
	addrs.Addr
	SourceInfo

	SchemaMigrationsTable string         `hcl:"schema_migrations_table"`
	MigrationsExpr        hcl.Expression `hcl:"migrations"`
	Migrations            []Migration    // populated from refs in expression
}

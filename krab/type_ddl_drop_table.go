package krab

// DDLDropTable contains DSL for dropping tables.
type DDLDropTable struct {
	Table string `hcl:"table,label"`
}

package krab

// DDLForeignKey constraint DSL for table DDL.
type DDLForeignKey struct {
	Name       string        `hcl:"name,label"`
	Columns    []string      `hcl:"columns"`
	References DDLReferences `hcl:"references,block"`
}

package krab

// DDLPrimaryKey constraint DSL for table DDL.
type DDLPrimaryKey struct {
	Name    string   `hcl:"name,label"`
	Columns []string `hcl:"columns"`
	Include []string `hcl:"include,optional"`
}

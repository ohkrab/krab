package krab

// DDLUnique constraint DSL for table DDL.
type DDLUnique struct {
	Name    string   `hcl:"name,label"`
	Columns []string `hcl:"columns"`
	Include []string `hcl:"include,optional"`
}

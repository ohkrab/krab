package krab

// DDLCreateTable contains DSL for creating tables.
type DDLCreateTable struct {
	Name        string           `hcl:"name,label"`
	Unlogged    bool             `hcl:"unlogged,optional"`
	Columns     []*DDLColumn     `hcl:"column,block"`
	PrimaryKeys []*DDLPrimaryKey `hcl:"primary_key,block"`
	ForeignKeys []*DDLForeignKey `hcl:"foreign_key,block"`
	Uniques     []*DDLUnique     `hcl:"unique,block"`
	Checks      []*DDLCheck      `hcl:"check,block"`
}

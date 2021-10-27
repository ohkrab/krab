package krab

// DDLReferences DSL for ForeignKey.
type DDLReferences struct {
	Table    string   `hcl:"table,label"`
	Columns  []string `hcl:"columns"`
	OnDelete string   `hcl:"on_delete,optional"`
	OnUpdate string   `hcl:"on_update,optional"`
}

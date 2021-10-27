package krab

// DDLCheck constraint DSL for table DDL.
type DDLCheck struct {
	Name       string `hcl:"name,label"`
	Expression string `hcl:"expression"`
}

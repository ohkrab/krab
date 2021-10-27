package krab

import "github.com/hashicorp/hcl/v2"

// DDLColumn constraint DSL for table DDL.
type DDLColumn struct {
	Name      string              `hcl:"name,label"`
	Type      string              `hcl:"type,label"`
	Null      *bool               `hcl:"null,optional"`
	Identity  *DDLIdentity        `hcl:"identity,block"`
	Default   hcl.Expression      `hcl:"default,optional"`
	Generated *DDLGeneratedColumn `hcl:"generated,block"`
}

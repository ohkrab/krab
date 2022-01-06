package krab

import (
	"io"

	"github.com/ohkrab/krab/krabhcl"
)

// Action represents custom action to execute.
//
type Action struct {
	Namespace string `hcl:"namespace,label"`
	RefName   string `hcl:"ref_name,label"`

	Arguments *Arguments `hcl:"arguments,block"`

	SQL string `hcl:"sql"`
}

func (a *Action) Addr() krabhcl.Addr {
	return krabhcl.Addr{Keyword: "action", Labels: []string{a.Namespace, a.RefName}}
}

func (a *Action) InitDefaults() {
	if a.Arguments == nil {
		a.Arguments = &Arguments{}
	}
	a.Arguments.InitDefaults()
}

func (a *Action) Validate() error {
	return ErrorCoalesce(
		ValidateRefName(a.Namespace),
		ValidateRefName(a.RefName),
	)
}

func (m *Action) ToSQL(w io.StringWriter) {
	w.WriteString(m.SQL)
}

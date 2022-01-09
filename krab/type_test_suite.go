package krab

import (
	"github.com/ohkrab/krab/krabhcl"
)

// TestSuite represents test runner configuration.
//
type TestSuite struct {
	RefName string          `hcl:"ref_name,label"`
	Before  TestSuiteBefore `hcl:"before,block"`

	Tests []*TestExample
}

func (t *TestSuite) Addr() krabhcl.Addr {
	return krabhcl.Addr{Keyword: "test_suite", Labels: []string{t.RefName}}
}

func (t *TestSuite) Validate() error {
	return ErrorCoalesce(
		ValidateRefName(t.RefName),
	)
}

type TestSuiteBefore struct {
	Dos []*Do `hcl:"do,block"`
}

package krab

import (
	"github.com/ohkrab/krab/krabhcl"
)

// TestSuite represents test runner configuration.
type TestSuite struct {
	Tests   []*TestExample
}

func (t *TestSuite) Addr() krabhcl.Addr {
	return krabhcl.NullAddr
}

func (t *TestSuite) Validate() error {
	return nil
}

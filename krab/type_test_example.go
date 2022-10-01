package krab

import (
	"fmt"
	"strconv"

	"github.com/ohkrab/krab/krabhcl"
)

// TestExample represents test runner configuration.
//
type TestExample struct {
	TestSuiteRefName string                `hcl:"test_suite,label"`
	Name             string                `hcl:"name,label"`
	Set              *SetRuntimeParameters `hcl:"set,block"`
	Its              []*TestExampleIt      `hcl:"it,block"`
}

func (t *TestExample) Addr() krabhcl.Addr {
	return krabhcl.Addr{Keyword: "test", Labels: []string{t.TestSuiteRefName, t.Name}}
}

func (t *TestExample) Validate() error {
	return ErrorCoalesce(
		ValidateRefName(t.TestSuiteRefName),
	)
}

// TestExampleIt represents one use case for test example that contain assertions.
type TestExampleIt struct {
	Comment     string        `hcl:"comment,label"`
	Do          *Do           `hcl:"do,block"`
	RowAsserts  []*AssertRow  `hcl:"row,block"`
	RowsAsserts []*AssertRows `hcl:"rows,block"`
}

// AssertRows
type AssertRows struct {
	Expectations []*Expect `hcl:"expect,block"`
}

// AssertRow
type AssertRow struct {
	Scope        string    `hcl:"scope,label"`
	Expectations []*Expect `hcl:"expect,block"`
}

// EachRow yields each function for every row defined in the scope.
func (a *AssertRow) EachRow(each func(i int64)) {
	i, err := strconv.ParseInt(a.Scope, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse scope `%s` to rows: %v", a.Scope, err))
	}
	each(i)
}

// Expect
type Expect struct {
	Subject  string  `hcl:"subject,label"`
	Equal    *string `hcl:"eq,optional"`
	Contains *string `hcl:"contains,optional"`
}

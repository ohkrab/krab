package spec

import (
	"testing"
)
func TestTest(t *testing.T) {
	t.Skip()
	return
	c := mockCli(mockConfig(`
test "spec name" {
  it "successful scenario" {
    query "SELECT 1 as c UNION ALL SELECT 2" {
      row "0" {
        col "must be 1" { assert = "c = 1" }
      }

      row "1" {
        col "must be 2" { assert = "c = 2" }
      }
    }
  }

  it "failed scenario" {
    query "SELECT 0 as c" {
      row "0" {
        col "must be 1" { assert = "c = 1" }
      }
    }
  }

  xit "skipped scenario" {
    query "SELECT 1 as c" {
      row "0" {
        col "must be 1" { assert = "c = 1" }
      }
    }
  }
}
`))
	defer c.Teardown()

	c.AssertSuccessfulRun(t, []string{"test"})
}

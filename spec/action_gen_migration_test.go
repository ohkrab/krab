package spec

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActionGenMigration(t *testing.T) {
	c := mockCli(mockConfig(``))
	defer c.Teardown()
	c.AssertSuccessfulRun(t, []string{"gen", "migration", "-name", "create_maps"})
	c.AssertOutputContains(t, "migration.create_maps")
	files := c.FSFiles()
	assert.Len(t, files, 1)
	for k, b := range files {
		expected := `migration "create_maps" {
  version = "20230101"

  up {
  }

  down {
  }
}`
		ok, err := c.fs.FileContainsBytes(k, []byte(expected))
		assert.NoError(t, err)
		if !ok {
			fmt.Println("Expected:", expected)
			fmt.Println("Current:", string(b))
			assert.FailNow(t, "Output file does not contain valid data")
		}
	}
}

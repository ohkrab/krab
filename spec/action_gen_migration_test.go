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

func TestActionGenMigrationWithParams(t *testing.T) {
	c := mockCli(mockConfig(``))
	defer c.Teardown()
	c.AssertSuccessfulRun(t, []string{
		"gen", "migration", "-name", "create_maps",
		"id", "name:varchar", "project_id:bigint:fk=projects,id", "timestamps",
	})
	c.AssertOutputContains(t, "migration.create_maps")
	files := c.FSFiles()
	assert.Len(t, files, 1)
	for k, b := range files {
		expected := `migration "create_maps" {
  version = "20230101"

  up {
    create_table "maps" {
      column "id" "bigint" {
        identity {}
      }
      column "name" "varchar" {}
      column "created_at" "timestamptz" {
		null = false
	  }
      column "updated_at" "timestamptz" {
		null = false
	  }
      primary_key {
		columns = ["id"]
	  }
      foreign_key {
        columns = ["project_id"]

        references "projects" {
          columns = ["id"]

          on_delete = "restrict"
          on_update = "cascade"
        }
      }
    }
  }

  down {
    drop_table "maps" {}
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

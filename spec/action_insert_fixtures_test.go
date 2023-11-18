package spec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActionInsertFixtures(t *testing.T) {
	c := mockCli(mockConfig(`
migration "create_fakes" {
  version = "v1"

  up   { sql = "CREATE TABLE fakes(name VARCHAR)" }
  down { sql = "DROP TABLE fakes" }
}

migration_set "public" {
  migrations = [migration.create_fakes]
}

action "seed" "fakes" {
  description = "Insert fake data into the fakes table"

  sql = <<-SQL
		INSERT INTO fakes(name) VALUES
		({{ fake "Food.Fruit" | quote }}),
		({{ fake "Internet.SafeEmail" | quote }}),
		({{ fake "Internet.Domain" | quote }}),
		({{ fake "Internet.Ipv4" | quote }}),
		({{ fake "Person.FirstName" | quote }}),
		({{ fake "Person.LastName" | quote }}),
		({{ fake "Person.Name" | quote }}),
		({{ fake "Address.CountryCode" | quote }}),
		({{ fake "Color.Hex" | quote }})
	SQL
}
`))
	defer c.Teardown()

	c.AssertSuccessfulRun(t, []string{"migrate", "up", "public"})
	c.AssertSchemaMigrationTable(t, "public", "v1")

	c.AssertSuccessfulRun(t, []string{"action", "seed", "fakes"})

	cols, rows := c.Query(t, "SELECT * FROM fakes")
	assert.ElementsMatch(t, []string{"name"}, cols, "Columns must match")
	if assert.Equal(t, 9, len(rows)) {
		for i := range rows {
			assert.NotEmpty(t, rows[i]["name"])
		}
	}
}

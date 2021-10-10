package krab

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParserArguments(t *testing.T) {
	assert := assert.New(t)

	p := mockParser(
		"src/public.krab.hcl",
		`
migration_set "public" {
  arguments {
    arg "str" {
	  type     = "string"
	}

	arg "num" {
	  type     = "int"
	}

	arg "none" {}
  }

  migrations = []
}
`)
	c, err := p.LoadConfigDir("src")
	assert.NoError(err, "Arguments should be parsed")

	set := c.MigrationSets["public"]
	assert.NotNil(set)
	assert.Equal(3, len(set.Arguments.Args))

	assert.Equal("str", set.Arguments.Args[0].Name)
	assert.Equal("string", set.Arguments.Args[0].Type)

	assert.Equal("num", set.Arguments.Args[1].Name)
	assert.Equal("int", set.Arguments.Args[1].Type)

	assert.Equal("none", set.Arguments.Args[2].Name)
	assert.Equal("string", set.Arguments.Args[2].Type)
}

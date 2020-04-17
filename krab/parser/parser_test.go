package parser

import (
	"testing"
)

func TestMigration(t *testing.T) {
	parsed, err := ParseFromFile("test/fixtures/migrations/create_table.hcl")

	expectNoError(t, err, func() {
		t.Log(parsed.Ast)
	})
}

func expectNoError(t *testing.T, err error, callback func()) {
	if err != nil {
		t.Error(err)
	} else {
		callback()
	}
}

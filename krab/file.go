package krab

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
)

// File represents all resource definitions within a single file.
type File struct {
	File *hcl.File

	Migrations    []*Migration
	MigrationSets []*MigrationSet
	Actions       []*Action
	TestSuites    []*TestSuite
	TestExamples  []*TestExample
	Wasms         []*WebAssembly
}

var schemaFile = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "migration",
			LabelNames: []string{"name"},
		},
		{
			Type:       "migration_set",
			LabelNames: []string{"name"},
		},
		{
			Type:       "action",
			LabelNames: []string{"namespace", "name"},
		},
	},
	Attributes: []hcl.AttributeSchema{},
}

// Decode parses HCL into struct.
func (f *File) Decode(ctx *hcl.EvalContext) error {
	f.Migrations = []*Migration{}
	f.MigrationSets = []*MigrationSet{}
	f.Actions = []*Action{}
	f.TestSuites = []*TestSuite{}
	f.TestExamples = []*TestExample{}
	f.Wasms = []*WebAssembly{}

	content, diags := f.File.Body.Content(schemaFile)
	if diags.HasErrors() {
		return fmt.Errorf("failed to decode file body: %s", diags.Error())
	}

	for _, b := range content.Blocks {
		switch b.Type {
		case "migration":
			migration := new(Migration)
			err := migration.DecodeHCL(ctx, b)
			if err != nil {
				return err
			}
			f.Migrations = append(f.Migrations, migration)
		}
		fmt.Println(b.Type)
	}

	return nil
}

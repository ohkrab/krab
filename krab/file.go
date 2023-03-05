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
	TestSuite     *TestSuite
	TestExamples  []*TestExample
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
		{
			Type:       "test",
			LabelNames: []string{"name"},
		},
	},
	Attributes: []hcl.AttributeSchema{},
}

// Decode parses HCL into struct.
func (f *File) Decode(ctx *hcl.EvalContext) error {
	f.Migrations = []*Migration{}
	f.MigrationSets = []*MigrationSet{}
	f.Actions = []*Action{}
	f.TestSuite = &TestSuite{}
	f.TestExamples = []*TestExample{}

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

		case "migration_set":
			migrationSet := new(MigrationSet)
			err := migrationSet.DecodeHCL(ctx, b)
			if err != nil {
				return err
			}
			f.MigrationSets = append(f.MigrationSets, migrationSet)

		case "action":
			action := new(Action)
			err := action.DecodeHCL(ctx, b)
			if err != nil {
				return err
			}
			f.Actions = append(f.Actions, action)

		case "test":
			test := new(TestExample)
			err := test.DecodeHCL(ctx, b)
			if err != nil {
				return err
			}
			f.TestExamples = append(f.TestExamples, test)

		default:
			return fmt.Errorf("Unknown block `%s`", b.Type)
		}
	}

	return nil
}

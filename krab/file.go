package krab

import "github.com/hashicorp/hcl/v2"

// File represents all resource definitions within a single file.
type File struct {
	Migrations    []*Migration    `hcl:"migration,block"`
	MigrationSets []*MigrationSet `hcl:"migration_set,block"`
	Actions       []*Action       `hcl:"action,block"`

	Raw *RawFile
}

// RawFile represents all resource definitions within a single file before mapping to structs.
type RawFile struct {
	Migrations []*RawMigration `hcl:"migration,block"`
	Remain     hcl.Body        `hcl:",remain"`
}

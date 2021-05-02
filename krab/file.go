package krab

// File represents all resource definication within a single file.
type File struct {
	Migrations    []*Migration    `hcl:"migration"`
	MigrationSets []*MigrationSet `hcl:"migration_set"`
}

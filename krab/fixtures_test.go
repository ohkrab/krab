package krab

func CreateMigrationSet(
	refName string,
	migrationData ...string,
) *MigrationSet {
	migrations := make([]*Migration, len(migrationData)/3)

	for i := 0; i < len(migrationData); i += 3 {
		migrations[i/3] = &Migration{
			Version: migrationData[i],
			Up: MigrationUp{
				SQL: migrationData[i+1],
			},
			Down: MigrationDown{
				SQL: migrationData[i+2],
			},
		}
	}
	set := &MigrationSet{RefName: refName, Migrations: migrations, Hooks: &Hooks{}}
	set.InitDefaults()
	return set
}

package spec

import (
	"testing"

	"github.com/jmoiron/sqlx"
)

func XTestActionMigrateUpDsl(t *testing.T) {
	withPg(t, func(db *sqlx.DB) {
		c := mockCli(mockConfig(`
migration "create_categories" {
  version = "v1"

  up {
    create_table "categories" {
	  column "id" "bigint" {}
	  column "name" "varchar" { null = false }
	}
  }

  down {
    drop_table "categories" {}
  }
}

migration "create_animals" {
  version = "v2"

  up {
	create_table "animals" {
	  unlogged = false

	  column "id" "bigint" {
		identity "always" {}
	  }

	  column "name" "varchar" { null = false }
	  
	  column "extinct" "boolean" {
	    null    = false
		default = false
	  }

	  column "weight_kg" "int" { null = false }

	  column "weight_g" "int" {
		generated {
		  as = "weight_kg * 1000"
		}
	  }

	  column "category_id" "bigint" {
	    null = false
	  }

	  unique "unique_name" {
		columns = ["name"]
		include = ["weight_kg"]
	  }

	  primary_key "pk" {
	    columns = ["id"]
		include = ["name"]
	  }

	  check "ensure_positive_weight" {
	    expression = "length(weight_kg) > 0"
	  }

	  foreign_key "fk" {
	    columns = ["category_id"]

		references "categories" {
		  columns = ["id"]

		  on_delete = "cascade"
		  on_update = "cascade"
		}
	  }
	}
  }

  down {
    drop_table "animals" {}
  }
}

migration_set "animals" {
  migrations = [
    migration.create_categories,
    migration.create_animals
  ]
}
`))
		c.AssertSuccessfulRun(t, []string{"migrate", "up", "animals"})
		c.AssertOutputContains(t,
			`
create_categories v1
create_animals v2
Done
`,
		)
		c.AssertSchemaMigrationTable(t, db, "public", "v1")
	})
}

// 	create_index "idx_uniq_name" {
// 	  unique  = true
// 	  columns = ["name"]
// 	  using   = "btree"
// 	  include = ["weight_kg"]
// 	  concurrently = false
// 	}

// 	create_index "idx_heavy_animals" {
// 	  columns = ["weight_kg"]
// 	  where   = "weight_kg > 5000"
// 	}

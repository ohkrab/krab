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
	  column "id" "bigint" { primary_key = true }
	  column "name" "varchar" { null = false }
	  
	  constraint "unique_name" {
		unique = ["name"]
	  }
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
	  column "id" "bigint" {
		identity = "always"
	  }

	  column "name" "varchar" { null = false }

	  column "extinct" "boolean" {
	    null    = false
		default = false
	  }

	  column "weight_kg" "int" { null = false }

	  column "category_id" "bigint" {
	    null = false
	  }

	  constraint "pk" {
		primary_key = ["id"]
	  }

	  constraint "ensure_positive_weight" {
	    check = "length(weight_kg) > 0"
	  }

	  constraint "fk" {
	    columns = ["category_id"]
		references "categories" {
		  columns = ["id"]
		  on_delete = "cascade"
		  on_update = "no action"
		}
	  }
	}

	create_index "idx_uniq_name" {
	  unique  = true
	  columns = ["name"]
	  using   = "btree"
	  include = ["weight_kg"]
	}

	create_index "idx_heavy_animals" {
	  columns = ["weight_kg"]
	  where   = "weight_kg > 5000"
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
create_users v1
Done
`,
		)
		c.AssertSchemaMigrationTable(t, db, "public", "v1")
	})
}

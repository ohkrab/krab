package spec

import (
	"testing"
)

func TestActionMigrateUpDsl(t *testing.T) {
	c := mockCli(mockConfig(`
migration "create_categories" {
  version = "v1"

  up {
    create_table "categories" {
	  column "id" "bigint" {}
	  column "name" "varchar" { null = false }

	  primary_key { columns = ["id"] }
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
	  unlogged = true

	  column "id" "bigint" {
		identity {}
	  }

	  column "name" "varchar" { null = true }
	  
	  column "extinct" "boolean" {
	    null    = false
		default = true
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

	  unique {
		columns = ["name"]
		include = ["weight_kg"]
	  }

	  primary_key {
	    columns = ["id"]
		include = ["name"]
	  }

	  check "ensure_positive_weight" {
	    expression = "weight_kg > 0"
	  }

	  foreign_key {
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
	defer c.Teardown()
	if c.AssertSuccessfulRun(t, []string{"migrate", "up", "animals"}) {
		c.AssertSchemaMigrationTable(t, "public", "v1", "v2")
		c.AssertOutputContains(t,
			`
create_categories v1
create_animals v2
Done
`,
		)
		c.AssertSQLContains(t, `
CREATE TABLE "categories"(
  "id" bigint,
  "name" varchar NOT NULL
, PRIMARY KEY ("id")
)
	`)
		c.AssertSQLContains(t, `
CREATE UNLOGGED TABLE "animals"(
  "id" bigint GENERATED ALWAYS AS IDENTITY,
  "name" varchar NULL,
  "extinct" boolean NOT NULL DEFAULT true,
  "weight_kg" int NOT NULL,
  "weight_g" int GENERATED ALWAYS AS (weight_kg * 1000) STORED,
  "category_id" bigint NOT NULL
, PRIMARY KEY ("id") INCLUDE ("name")
, FOREIGN KEY ("category_id") REFERENCES "categories"("id") ON DELETE cascade ON UPDATE cascade
, UNIQUE ("name") INCLUDE ("weight_kg")
, CONSTRAINT "ensure_positive_weight" CHECK (weight_kg > 0)
)
	`)

		if c.AssertSuccessfulRun(t, []string{"migrate", "down", "animals", "-version", "v2"}) {
			c.AssertSchemaMigrationTable(t, "public", "v1")
			c.AssertOutputContains(t,
				`
Done
`,
			)
			c.AssertSQLContains(t, `
DROP TABLE "animals"
	`)
		}
	}
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

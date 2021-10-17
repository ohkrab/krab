migration "add_tenants" {
  version = "202006_01"

  up {
    sql = "CREATE TABLE tenants(name varchar PRIMARY KEY)"
  }

  down {
    sql = "DROP TABLE tenants"
  }
}

migration "add_tenants_index" {
  version     = "202108_01"
  transaction = false

  up {
    sql = "CREATE INDEX CONCURRENTLY idx_tenants_name ON tenants(name)"
  }

  down {
    sql = "DROP INDEX CONCURRENTLY idx_tenants_name"
  }
}

migration_set "public" {
  migrations = [
    migration.add_tenants,
    migration.add_tenants_index,
  ]
}

migration "add_users" {
  version = "202006_01"

  up {
    sql = "CREATE TABLE users(email varchar PRIMARY KEY)"
  }

  down {
    sql = "DROP TABLE users"
  }
}

migration_set "tenant" {
  arguments {
    arg "schema" {
      type = "string"
    }
  }

  hooks {
    before = "SET search_path TO {{args.schema}}"
  }

  migrations = [
    migration.create_assets
  ]
}

migration "create_assets" {
  version = "202010_01"

  up {
    # create_table "assets" {
    #   column "id" {
    #     type = "int"
    #     identity {
    #       always = true
    #     }
    #   }

    #   column "name" {
    #     type = "varchar"
    #     null = false
    #   }

    #   column "size" {
    #     type    = "int"
    #     default = 0
    #   }


    #   constraint "ensure_positive_size" {
    #     check = "size > 0"
    #   }
    # }

    # alter_table "users" {
    #     add_column "email" {
    #       type = "varchar"
    #       null = true
    #     }

    #     drop_column "deprecated_field" {}

    #     primary_key = ["email", "name"]

    #     create_index "idx_uniq_emails" {
    #       unique  = true
    #       columns = ["email"]
    #     }

    #     constraint "users_pk" {
    #       columns = ["email"]
    #       check "valid_nu
    #     }
    # }
  }

  # down = up.reverse
  down {}

  # hooks {
  #     after "up" {
  #         do = wasm.file("../wasm/migrate_from_old_system.wasm")
  #     }
  # }
}

migration "add_tenants" {
  up {
    sql = "CREATE TABLE tenants(name varchar PRIMARY KEY)"
  }

  down {
    sql = "DROP TABLE tenants"
  }
}

migration_set "public" {
  schema_migrations_table = "krab_migrations"

  migrations = [
    migration.add_tenants
  ]
}

migration "add_users" {
  up {
    sql = "CREATE TABLE users(email varchar PRIMARY KEY)"
  }

  down {
    sql = "DROP TABLE users"
  }
}

migration_set "tenant" {
  schema_migrations_table = "krab_migrations"

  migrations = [
    migration.add_users
  ]
}

migration "create_users" {
  up {
    alter_table "users" {
        add_column "email" {
          type = "varchar"
          null = true
        }

        drop_column "deprecated_field" {}

        primary_key = ["email", "name"]

        create_index "idx_uniq_emails" {
          unique  = true
          columns = ["email"]
        }

        constraint "users_pk" {
          columns = ["email"]
          check "valid_nu
        }
    }
  }

  down = up.reverse

  hooks {
      after "up" {
          do = wasm.file("../wasm/migrate_from_old_system.wasm")
      }
  }
}

connection "pg" "test" {
  uri = "postgres://localhost:5432/krab_test"
}
# bind "pg_migration" "up" {
#   cli   = ["migrate"]
#   agent = true
# }

# # allow multiple source to migrate
# command "pg_migrate_up" "default" {
#   migrate_up {
#     # connection_uri = envs.DATABASE_URI
#     # sets = [pg_migration_set.public]
#   }

#   migrate_up {
#     # connection_uri = pg_connection.default
#     # sets = [pg_migration_set.tenant]
#   }

#   migrate_up {
#     # connection_uri = params.get.database_uri
#     # sets           = [pg_migration_set.tenant]
#   }
# }

# bind "pg_migration" "rollback" {
#   cli = ["db", "rollback"]

#   args = {
#     step = {
#       default = 1
#       type    = "number"
#     }
#   }

#   triggers = [
#     # command.pg_migrations_rollback.default
#   ]
# }

# command "pg_migrations_rollback" "default" {
#   input = {
#     # step = args.step
#   }

#   rollback {
#     # connection = pg_connection.default
#     # sets       = [pg_migration_set.public]
#   }
# }

resource "pg_connection" "default" {
  # uri = envs.DATABASE_URI
  # = vault.app.config.db_uri
  # = param.database_uri?
}


# resource "pg_migration" "add_tenants" {
#   up {
#     sql = <<SQL
#             CREATE TABLE ...
#         SQL
#   }

#   down {
#     sql = "..."
#   }
# }

# resource "pg_migration_set" "public" {
#   schema_info = "_migrations"

#   migrations = [
#     pg_migration.add_tenants.id
#   ]
# }
# resource "pg_migration_set" "tenant" {
#   schema_info = "_migrations"

#   migrations = [
#     # pg_migration.create_scans,
#     # pg_migration.create_refs,
#   ]
# }

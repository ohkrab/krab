# define plugins and its sources
# for now it will only support local files
# 
# krab init - will fetch plugin (for local file it will load it from its destination)
plugin "sshconfig" {
    file = "plugins/sshconfig"
}

plugin "pg_migration" {
    file = "plugins/pg_migration"
}

# register command
command "sshconfig_cmd" "generate_config" {
    # krab plugin ssh config
    cli = [
        "ssh",
        "config"
    ]

    params = {
        hostalias = string
        port = number
        user = string
        hostname = string
    }
}



# allow multiple source to migrate
command "pg_migration_up_cmd" "default" {
    cli = ["migrate"]

    migration {
        connection_uri = env.DATABASE_URI
        sets = [pg_migration_set.public]
    }

    migration {
        connection_uri = pg_connection.default
        sets = [pg_migration_set.tenant]
    }
}

command "pg_migration_up_cmd" "tenant" {
    cli = ["tenant", "migrate"]
    api = true # api will expose /plugin/tenant/migrate JSON api

    params = [
        database_uri = {
            type = string
            required = true
        }
    ]

    # commands can have params to be defined
    # so you either send them from CLI or JSON api
    migration {
        connection_uri = param.database_uri
        sets = [pg_migration_set.tenant]
    }
}

command "pg_migration_down_cmd" "default" {
    cli = ["migrate", "down"]

    migration {
        connection = pg_connection.default
        sets = [pg_migration_set.public]
    }

    migration {
        connection = pg_connection.default
        sets = [pg_migration_set.tenant]
    }
}

resource "pg_connection" "default" {
    uri = env.DATABASE_URI
    # = vault.app.config.db_uri
    # = param.database_uri?
}

resource "pg_migration_set" "public" {
    schema_info = "metadata"

    migrations = [
        pg_migration.add_tenants
    ]
}

resource "pg_migration" "add_tenants" {
    up {}
    down {}
}

resource "pg_migration_set" "tenant" {
    schema_info = "metadata"

    migrations = [
        pg_migration.create_scans,
        pg_migration.create_refs,
    ]
}

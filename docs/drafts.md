# drafts

doc for ideas.

## echo

```
plugin "echo" {
    file = "plugins/echo"
}

bind "echo" "print" {

    cli = ["echo"] # krab plugin echo

    params {
        name = string
    }
}


resource "echo" "print_hello_world" {
    in = bind.echo.print

    echo {
        stdout = "hello {{param.get.name}}!"
    }

    echo {
        stdout = "\n"
    }
}


bind "echo" "gen_sshconfig" {
    enabled = env.exists.ENABLE_SSH_TOOLS

    cli = ["sshconfig", "gen"]
    agent = true # availiabne in json api /plugin/ssconfig/gen

    # cli:
    #   krab plugin sshconfig gen -port 22 ...
    #
    # agent:
    #   http POST :8888/plugin/sshconfig/gen
    #   body: {"port": 22, ...}
    #
    params {
        hostalias = string
        port = number
        user = string
        hostname = string
    }
}

resource "echo" "gen_sshconfig" {
    in = bind.echo.generate_config

    echo {
         stdout = "host {{params.get.hostalias}}"
    }
    
    echo {
        stdout = <<EOS
Host {{params.get.hostalias}}
    HostName {{params.get.host}}
    Port {{params.get.port}}
    User {{parans.get.user}}
EOS
    }
    
# Host mars
#    HostName 192.168.1.1
#    Port 22
#    User elon
}

# define plugins and its sources
# for now it will only support local files
# 
# krab init - will fetch plugin (for local file it will load it from its destination)

plugin "pg_migration" {
    file = "plugins/pg_migration"
}

bind "pg_migration" "up" {
    cli = ["migrate"]
    agent = true
}

# allow multiple source to migrate
command "pg_migrate_up" "default" {
    migrate_up {
        connection_uri = envs.get.DATABASE_URI
        sets = [pg_migration_set.public]
    }

    migrate_up {
        connection_uri = pg_connection.default
        sets = [pg_migration_set.tenant]
    }

    migrate_up {
        connection_uri = params.get.database_uri
        sets = [pg_migration_set.tenant]
    }
}

bind "pg_migration" "rollback" {
    cli = ["db", "rollback"]

    args = {
        step = {
            default = 1
            type = number
        }
    }

    triggers = [
        command.pg_migrations_rollback.default
    ]
}

command "pg_migrations_rollback" "default" {
    input = {
        step = args.step
    }

    rollback {
        connection = pg_connection.default
        sets = [pg_migration_set.public]
    }
...
}

resource "pg_connection" "default" {
    uri = env.DATABASE_URI
    # = vault.app.config.db_uri
    # = param.database_uri?
}

resource "pg_migration_set" "public" {
    schema_info = "_migrations"

    migrations = [
        pg_migration.add_tenants
    ]
}

resource "pg_migration" "add_tenants" {
    up {
        sql = <<SQL
            CREATE TABLE ...
        SQL
    }

    down {
        sql = "..."
    }
}

resource "pg_migration_set" "tenant" {
    schema_info = "_migrations"

    migrations = [
        pg_migration.create_scans,
        pg_migration.create_refs,
    ]
}
```

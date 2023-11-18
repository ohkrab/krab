action "db" "create" {
  transaction = false
  description = "Create a database and assign database owner"

  arguments {
    arg "name" {
      description = "Database name"
    }

    arg "user" {
      description = "Database user"
    }
  }

  sql = "CREATE DATABASE {{`{{ .Args.name | quote_ident }}`}} OWNER {{`{{ .Args.user | quote_ident }}`}}"
}

action "user" "create" {
  description = "Create a database user with password"

  arguments {
    arg "user" {
      description = "Database user"
    }

    arg "password" {
      description = "Database password"
    }
  }

  sql = "CREATE USER {{`{{ .Args.user | quote_ident }}`}} WITH PASSWORD {{`{{ .Args.password | quote }}`}}"
}

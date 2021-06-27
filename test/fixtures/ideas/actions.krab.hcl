action {
  cli      = "mig:up"
  api_path = "mig/up"

  do {
    connection = connection.test

    migrate "up" {
      sets = [migration_set.public]
    }
  }
}

action {
  cli = "mig:all"

  do {
    connection = connection.test

    migrate "up" {
      sets = [migration_set.public]
    }
  }

  do {
    connection = connection.test

    for_each {
      in = query.select_tenants
    }

    migrate "up" {
      before {
          sql = "SET SESSION search_path TO ${pg.quote_ident(each.name)}"
      }

      sets = [migration_set.tenant]
    }
  }
}

action {
  cli = ["mig:down"]
  api = ["mig/down"]

  args {
    version = {
      type = "string"
    }
  }

  do {
    connection = connection.test

    migrate "down" {
      input {
        version = args.version
      }

      sets = [migration_set.public]
    }
  }
}

action {
  cli = ["mig:rollback"]
  api = ["mig/rollback"]

  args {
    step = {
      type     = "int"
      optional = true
      default  = 1
    }
  }

  do {
    connection = connection.test

    migrate "rollback" {
      input {
        step = args.step
      }

      sets = [migration_set.public]
    }
  }
}

migration "create_animals" {
  version = "v1"

  up { sql = "CREATE TABLE animals(name VARCHAR)" }
  down { sql = "DROP TABLE animals" }
}

migration "create_animals_view" {
  version = "v2"

  up { sql = "CREATE MATERIALIZED VIEW anims AS SELECT name FROM animals" }
  down { sql = "DROP MATERIALIZED VIEW anims" }
}

migration "seed_animals" {
  version = "v3"

  up { sql = "INSERT INTO animals(name) VALUES('Elephant'),('Turtle'),('Cat')" }
  down { sql = "TRUNCATE animals" }
}

migration_set "animals" {
  arguments {
    arg "name" {
      description = "Materialized view to be refreshed"
    }
  }

  migrations = [
    migration.create_animals,
    migration.create_animals_view,
    migration.seed_animals,
  ]
}

action "view" "refresh" {
  description = "Refresh a materialized view"
  
  arguments {
    arg "name" {
      description = "Materialized view to be refreshed"
    }
  }

  sql = "REFRESH MATERIALIZED VIEW {{ .Args.name | quote_ident }}"
}

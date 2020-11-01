globals {
  uri = "global_var"
}

connection "default" {
  uri = "postgres://krab:secret@localhost:5432/postgres"
}

connection "interpolated" {
  uri = "postgres://${env("USER")}:${env("PASSWORD")}@localhost:5432/postgres"
}

connection "referenced" {
  uri = global.uri
}

connection "duplicated" {
  uri = connection.default.uri
}

connection "from_env" {
  uri = env("PG_URI")
}


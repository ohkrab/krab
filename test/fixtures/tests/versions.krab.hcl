migration "create_version_type" {
  version = "v1"

  up {
    sql = <<SQL
      CREATE TYPE version AS (
        major SMALLINT,
        minor SMALLINT,
        patch SMALLINT
      )
    SQL
  }

  down { sql = "DROP TYPE version" }
}

migration "create_version_function" {
  version = "v2"

  up {
    sql = <<SQL
      CREATE FUNCTION version_inc(_ver version, _type varchar = 'major') RETURNS version
      AS
      $$
      DECLARE
        _v version;
      BEGIN
        _v := _ver;

        CASE 
          WHEN _type = 'major' THEN
            _v.major = _v.major + 1;
            _v.minor = 0;
            _v.patch = 0;
            RETURN _v;
        END CASE;

        RAISE EXCEPTION 'Failed to increase version using type = %s for version %L.%L.%L', _type, _ver.major, _ver.minor, _ver.patch;
      END;
      $$
      RETURNS NULL ON NULL INPUT
      LANGUAGE plpgsql
    SQL
  }

  down {
    sql = "DROP FUNCTION version_inc(version, varchar)"
  }
}

migration_set "versions" {
  arguments {
    arg "schema" {}
  }

  schema = "{{.Args.schema}}"

  migrations = [
    migration.create_version_type,
    migration.create_version_function
  ]
}

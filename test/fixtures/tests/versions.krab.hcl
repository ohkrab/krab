migration "create_version_type" {
  version = "v1"

  up {
    sql = <<SQL
      CREATE TYPE sem_version AS (
        major SMALLINT,
        minor SMALLINT,
        patch SMALLINT
      );
    SQL
  }

  down { sql = "DROP TYPE sem_version" }
}

migration "create_version_function" {
  version = "v2"

  up {
    sql = <<SQL
      CREATE FUNCTION version_inc(_ver sem_version, _type varchar = 'major') RETURNS sem_version
      AS
      $$
      DECLARE
        _v sem_version;
      BEGIN
        _v := _ver;

        CASE 
          WHEN _type = 'major' THEN
            _v.major = _v.major + 1;
            _v.minor = 0;
            _v.patch = 0;

          WHEN _type = 'minor' THEN
            _v.minor = _v.minor + 1;
            _v.patch = 0;

          WHEN _type = 'patch' THEN
            _v.patch = _v.patch + 1;

          ELSE
            RAISE EXCEPTION 'Failed to increase version using type = % for version %.%.%', _type, _ver.major, _ver.minor, _ver.patch;
        END CASE;

        RETURN _v;
      END;
      $$
      RETURNS NULL ON NULL INPUT
      LANGUAGE plpgsql
    SQL
  }

  down {
    sql = "DROP FUNCTION version_inc(sem_version, varchar)"
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

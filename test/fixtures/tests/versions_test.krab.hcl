test_suite "versions" {
  before {
    do {
      inputs = {
        schema = "testing"
      }

      migrate "up" { migration_set = migration_set.versions }
      migrate "down" {
        migration_set = migration_set.versions
        inputs        = { version = "v2" }
      }
      migrate "up" { migration_set = migration_set.versions }
    }
  }
}

test "versions" "version_inc()" {
  set {
    search_path = "testing"
  }

  it "increases `major` component by default when no type specified" {
    do {
      sql = "SELECT version_inc(row(1,1,1)::sem_version) AS ver"
    }

    row "0" {
      column "ver" { assert = "(2,0,0)" }
    }
  }

  it "increases `major` component and resets `minor` and `patch`" {
    do { sql = "SELECT version_inc(row(1,1,1)::sem_version, 'major') AS ver" }

    row "0" {
      column "ver" { assert = "ver = row(2,0,0)::sem_version" }
    }
  }

  describe "increases `minor` component and resets `patch` leaving `major` untouched" {
    do { sql = "SELECT version_inc(row(1,1,1)::sem_version, 'minor') AS ver" }

    # v1 - set scope
    row "0" {
      it "ver" { expect = "ver = row(2,0,0)::sem_version" }
      its { expect = "ver = row(2,0,0)::sem_version" }
    }

    # v2
    it "0" "ver" { expect = "ver = row(2,0,0)::sem_version" }
    its "0" { expect = "ver = row(2,0,0)::sem_version" }
  }

  it "increases `patch` component and leaves `major` and `minor` untouched" {
    do { sql = "SELECT version_inc(row(1,1,1)::sem_version, 'patch') AS ver" }

    row "0" {
      column "ver" { assert = "(1,1,2)" }
    }
  }

  /* it "raises error when increasing invalid component" { */
  /*   do { sql = "SELECT version_inc(row(1,1,1)::sem_version, 'invalid') AS ver" } */

  /*   rows { */
  /*     raises_error { */
  /*       message = "Failed to increase version using type = 'invalid' for version 1.1.1" */
  /*     } */
  /*   } */
  /* } */

  it "returns null on null input" {
    do { sql = "SELECT version_inc(null) AS ver" }

    row "0" {
      column "ver" { assert = null }
    }
  }
}

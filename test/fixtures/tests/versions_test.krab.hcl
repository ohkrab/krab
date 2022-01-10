test_suite "versions" {
  before {
    do {
      migration_set = migration_set.versions
      inputs = {
        schema = "aaa"
      }
    }
  }
}

test "versions" "version_inc()" {
  it "increases `major` component by default when no type specified" {
    do { sql = "SELECT version_inc('1.1.1') AS ver" }

    row "0" {
      expect "ver" { equal = "2.0.0" }
    }
  }

  it "increases `major` component and resets `minor` and `patch`" {
    do { sql = "SELECT version_inc('1.1.1', 'major') AS ver" }

    row "0" {
      expect "ver" { equal = "2.0.0" }
    }
  }

  it "increases `minor` component and resets `patch` leaving `major` untouched" {
    do { sql = "SELECT version_inc('1.1.1', 'minor') AS ver" }

    row "0" {
      expect "ver" { equal = "1.2.0" }
    }
  }

  it "increases `patch` component and leaves `major` and `minor` untouched" {
    do { sql = "SELECT version_inc('1.1.1', 'patch') AS ver" }

    row "0" {
      expect "ver" { equal = "1.1.2" }
    }
  }

  it "raises error when increasing invalid component" {
    do { sql = "SELECT version_inc('1.1.1', 'invalid') AS ver" }

    rows {
      expect "error" {
        equal = "Failed to increase version using type = 'invalid' for version 1.1.1"
      }
    }
  }

  it "returns null on null input" {
    do { sql = "SELECT version_inc(null) AS ver" }

    row "0" {
      expect "ver" { equal = null }
    }
  }
}

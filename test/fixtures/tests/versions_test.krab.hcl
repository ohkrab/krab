test "versions" "version_inc()" {
  it "increases `major` component by default when no type specified" {
    do {
      sql = "SELECT version_inc(row(1,1,1)::sem_version) AS ver"
    }

    row "0" {
      col "must be 2.0.0" { assert = "ver = row(2,0,0)::sem_version" }
    }
  }

  it "increases `major` component and resets `minor` and `patch`" {
    do { sql = "SELECT version_inc(row(1,1,1)::sem_version, 'major') AS ver" }

    row "0" {
      col "must be 2.0.0" { assert = "ver = row(2,0,0)::sem_version" }
    }
  }

  it "increases `minor` component and resets `patch` leaving `major` untouched" {
    do { sql = "SELECT version_inc(row(1,1,1)::sem_version, 'minor') AS ver" }

    # v1 - set scope
    row "0" {
      col "must be 1.2.0" { assert = "ver = row(1,2,0)::sem_version" }
    }
  }

  it "increases `patch` component and leaves `major` and `minor` untouched" {
    do { sql = "SELECT version_inc(row(1,1,1)::sem_version, 'patch') AS ver" }

    row "0" {
      col "must be 1.1.2" { assert = "ver = row(1,1,2)::sem_version" }
    }
  }

  it "returns null on null input" {
    do { sql = "SELECT version_inc(null) AS ver" }

    row "0" {
      col "must be null" { assert = "ver IS NULL" }
    }
  }
}

test "versions" "version_inc()" {
  query "SELECT version_inc(row(1,1,1)::sem_version) AS ver" {
    row "0" {
      it "increases major version, resets others" { assert = "ver = row(2,0,0)::sem_version" }
    }
  }

  query "SELECT version_inc(row(1,1,1)::sem_version, 'major') AS ver" {
    row "0" {
      it "increases major version, resets others" { assert = "ver = row(2,0,0)::sem_version" }
    }
  }

  query "SELECT version_inc(row(1,1,1)::sem_version, 'minor') AS ver" {
    row "0" {
      it "increases minor, resets patch" { assert = "ver = row(1,2,0)::sem_version" }
    }
  }

  query "SELECT version_inc(row(1,1,1)::sem_version, 'patch') AS ver" {
    row "0" {
      it "increases patch" { assert = "ver = row(1,1,2)::sem_version" }
    }
  }

  query "SELECT version_inc(null) AS ver" {
    row "0" {
      it "returns null on null input" { assert = "ver IS NULL" }
    }
  }
}

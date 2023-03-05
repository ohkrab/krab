test "version_inc()" {
  it "tests version changes" {
    query "SELECT version_inc(row(1,1,1)::sem_version) AS ver" {
      row "0" {
        col "increases major version, resets others" { assert = "ver = row(2,0,0)::sem_version" }
      }
    }

    query "SELECT version_inc(row(1,1,1)::sem_version, 'major') AS ver" {
      row "0" {
        col "increases major version, resets others" { assert = "ver = row(2,0,0)::sem_version" }
      }
    }

    query "SELECT version_inc(row(1,1,1)::sem_version, 'minor') AS ver" {
      row "0" {
        col "increases minor, resets patch" { assert = "ver = row(1,2,0)::sem_version" }
      }
    }

    query "SELECT version_inc(row(1,1,1)::sem_version, 'patch') AS ver" {
      row "0" {
        col "increases patch" { assert = "ver = row(1,1,2)::sem_version" }
      }
    }

    query "SELECT version_inc(null) AS ver" {
      row "0" {
        col "returns null on null input" { assert = "ver IS NULL" }
      }
    }
  }

  xit "test version on strings" {
    query "SELECT version_inc('1.1.1') AS ver" {
    }
  }
}

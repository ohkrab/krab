package krab

import (
	"strings"
	"testing"

	"github.com/franela/goblin"
)

func TestParser(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("Simple migration resource", func() {
		g.It("Should parse config without errors", func() {
			p := mockParser(
				"src/public.krab.hcl",
				`
migration "create_tenants" {
  up {
	sql = "CREATE TABLE tenants(name VARCHAR PRIMARY KEY)"
  }

  down {
	sql = "DROP TABLE tenants"
  }
}
`)
			c, err := p.LoadConfigDir("src")
			g.Assert(err).IsNil()

			if migration, ok := c.Migrations["create_tenants"]; ok {
				g.Assert(migration.RefName).Eql("create_tenants")
				g.Assert(migration.Up.Sql).Eql("CREATE TABLE tenants(name VARCHAR PRIMARY KEY)")
				g.Assert(migration.Down.Sql).Eql("DROP TABLE tenants")
			} else {
				g.Failf("Can't get migration %s", "create_tenants")
			}

		})
	})

	g.Describe("Optional content in up/down blocks for migrations", func() {
		g.It("Parses successfuly without providing up/down details", func() {
			p := mockParser(
				"src/public.krab.hcl",
				`migration "abc" {
                  up {}
				  down {}
				}`,
			)
			c, err := p.LoadConfigDir("src")
			g.Assert(err).IsNil()
			if migration, ok := c.Migrations["abc"]; ok {
				g.Assert(migration.RefName).Eql("abc")
				g.Assert(migration.Up.Sql).Eql("")
				g.Assert(migration.Down.Sql).Eql("")
			} else {
				g.Failf("Can't get migration %s", "abc")
			}
		})
	})

	g.Describe("Duplicated migration resource with the same ref name", func() {
		g.It("Config parsing should fail because of duplicates", func() {
			p := mockParser(
				"src/public.krab.hcl",
				`
migration "abc" {
  up { sql = "" }
  down { sql = "" }
}

migration "abc" {
  up { sql = "" }
  down { sql = "" }
}
`)
			_, err := p.LoadConfigDir("src")
			g.Assert(err).IsNotNil()
			g.Assert(strings.Contains(err.Error(), "Migration with the name 'abc' already exists")).IsTrue("Names must be unique")
		})
	})

	g.Describe("Simple migration_set resource", func() {
		g.It("Should parse config without errors", func() {
			p := mockParser(
				"src/migrations.krab.hcl",
				`
migration "abc" {
  up {}
  down {}
}

migration "def" {
  up {}
  down {}
}

migration "xyz" {
  up {}
  down {}
}
`,
				"src/sets.krab.hcl",
				`
migration_set "public" {
  migrations = [
  	migration.abc,
	migration.def,
  ]
}

migration_set "private" {
  migrations = [migration.xyz]
}
`)
			c, err := p.LoadConfigDir("src")
			g.Assert(err).IsNil()

			// public set
			publicSet, ok := c.MigrationSets["public"]
			if !ok {
				g.Fail("Failed to fetch 'public' set")
			}
			g.Assert(publicSet.RefName).Eql("public")
			g.Assert(len(publicSet.Migrations)).Eql(2)
			g.Assert(publicSet.Migrations[0].RefName).Eql("abc")
			g.Assert(publicSet.Migrations[1].RefName).Eql("def")

			// private set
			privateSet, ok := c.MigrationSets["private"]
			if !ok {
				g.Fail("Failed to fetch 'private' set")
			}
			g.Assert(privateSet.RefName).Eql("private")
			g.Assert(len(privateSet.Migrations)).Eql(1)
			g.Assert(privateSet.Migrations[0].RefName).Eql("xyz")
		})
	})

	g.Describe("Duplicated migration_set resource with the same ref name", func() {
		g.It("Config parsing should fail because of duplicates", func() {
			p := mockParser(
				"src/sets.krab.hcl",
				`
migration_set "abc" {
  migrations = []
}

migration_set "abc" {
  migrations = []
}
`)
			_, err := p.LoadConfigDir("src")
			g.Assert(err).IsNotNil()
			g.Assert(strings.Contains(err.Error(), "Migration Set with the name 'abc' already exists")).IsTrue("Names must be unique")
		})
	})
}

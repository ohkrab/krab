package krab

import (
	"strings"
	"testing"

	"github.com/franela/goblin"
	"github.com/spf13/afero"
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
			g.Assert(strings.Contains(err.Error(), "Migration with the name 'abc' already exists")).IsTrue()
		})
	})
}

// mockParser expects args: "path", "content", "path2", "content2", ...
func mockParser(pathContentPair ...string) *Parser {
	memfs := afero.NewMemMapFs()

	for i := 1; i < len(pathContentPair); i += 2 {
		path := pathContentPair[i-1]
		content := pathContentPair[i]
		afero.WriteFile(
			memfs,
			path,
			[]byte(content),
			0644,
		)
	}

	p := NewParser()
	p.fs = afero.Afero{Fs: memfs}
	return p
}

package krab

import (
	"testing"

	"github.com/franela/goblin"
	"github.com/spf13/afero"
)

func TestParser(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("Simple migrations with sets", func() {
		g.It("Should parse config without errors", func() {
			p := mockParser("src/a.krab.hcl", `
`)
			_, err := p.LoadConfigDir("src")
			g.Assert(err).IsNotNil("Parsing src/ should return error")

		})
		// 1. load dir
		// 2. parse with no errors
		// 3. verify correct  reference
	})
}

// mockParser expects args: "path", "content", "path2", "content2", ...
func mockParser(pathContentPair ...string) *Parser {
	memfs := afero.NewMemMapFs()

	for i := 1; i < len(pathContentPair); i += 2 {
		afero.WriteFile(
			memfs,
			pathContentPair[i-1],
			[]byte(pathContentPair[i]),
			0644,
		)
	}

	p := NewParser()
	p.fs = afero.Afero{Fs: memfs}
	return p
}

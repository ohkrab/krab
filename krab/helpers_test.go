package krab

import (
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/spf13/afero"
)

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
	p.FS = afero.Afero{Fs: memfs}
	return p
}

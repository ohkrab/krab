package krab

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/ohkrab/krab/krabhcl"
	"github.com/spf13/afero"
)

// CmdGenMigration generates migation file.
type CmdGenMigration struct {
	FS afero.Afero
}

// ResponseGenMigration json
type ResponseGenMigration struct {
	Path string `json:"path"`
	Ref  string `json:"ref"`
}

func (c *CmdGenMigration) Addr() krabhcl.Addr { return krabhcl.NullAddr }

func (c *CmdGenMigration) Name() []string {
	return append([]string{"gen", "migration"})
}

func (c *CmdGenMigration) HttpMethod() string { return "" }

func (c *CmdGenMigration) Do(ctx context.Context, o CmdOpts) (interface{}, error) {

	return c.run(ctx, args)
}

func (c *CmdGenMigration) run(ctx context.Context) (ResponseGenMigration, error) {
	result := ResponseGenMigration{}
	buf := bytes.Buffer{}

	ref := "create_animals"
	result.Ref = fmt.Sprint("migration.", ref)
	version := time.Now().UTC().Format("20060102_150405") // YYYYMMDD_HHMMSS

	buf.WriteString(`migration "`)
	buf.WriteString(ref)
	buf.WriteString(`" {`)
	buf.WriteString("\n")
	buf.WriteString(`  version = "`)
	buf.WriteString(version)
	buf.WriteString(`"`)
	buf.WriteString("\n\n")
	buf.WriteString(`  up {`)
	buf.WriteString("\n")
	buf.WriteString(`  }`)
	buf.WriteString("\n\n")
	buf.WriteString(`  down {`)
	buf.WriteString("\n")
	buf.WriteString(`  }`)
	buf.WriteString("\n")
	buf.WriteString(`}`)
	buf.WriteString("\n")

	c.FS.WriteFile("/tmp/migrate.krab.hcl", buf.Bytes(), 0644)

	return result, nil
}

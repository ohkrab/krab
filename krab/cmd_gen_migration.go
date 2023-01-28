package krab

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"

	"github.com/ohkrab/krab/krabenv"
	"github.com/ohkrab/krab/krabhcl"
	"github.com/spf13/afero"
)

// CmdGenMigration generates migation file.
type CmdGenMigration struct {
	FS afero.Afero
	VersionGenerator
}

// ResponseGenMigration json
type ResponseGenMigration struct {
	Path string `json:"path"`
	Ref  string `json:"ref"`
}

func (c *CmdGenMigration) Arguments() *Arguments {
	return &Arguments{
		Args: []*Argument{
			{
				Name:        "name",
				Type:        "string",
				Description: "Migration name",
			},
		},
	}
}

func (c *CmdGenMigration) Addr() krabhcl.Addr { return krabhcl.NullAddr }

func (c *CmdGenMigration) Name() []string {
	return append([]string{"gen", "migration"})
}

func (c *CmdGenMigration) HttpMethod() string { return "" }

func (c *CmdGenMigration) Do(ctx context.Context, o CmdOpts) (interface{}, error) {
	err := c.Arguments().Validate(o.Inputs)
	if err != nil {
		return nil, err
	}
	return c.run(ctx, o.Inputs)
}

func (c *CmdGenMigration) run(ctx context.Context, inputs Inputs) (ResponseGenMigration, error) {
	result := ResponseGenMigration{}

	dir, err := krabenv.ConfigDir()
	if err != nil {
		return result, err
	}
	dir = filepath.Join(dir, "db", "migrations")
	err = c.FS.MkdirAll(dir, 0755)
	if err != nil {
		return result, err
	}

	version := c.VersionGenerator.Next()
	ref := inputs["name"].(string)
	result.Ref = fmt.Sprint("migration.", ref)
	result.Path = filepath.Join(dir, fmt.Sprint(version, "_", ref, krabenv.Ext()))

	buf := bytes.Buffer{}
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

	c.FS.WriteFile(result.Path, buf.Bytes(), 0644)

	return result, nil
}

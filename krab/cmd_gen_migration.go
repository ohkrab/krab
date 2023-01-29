package krab

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ohkrab/krab/krabenv"
	"github.com/ohkrab/krab/krabhcl"
	"github.com/spf13/afero"
)

// CmdGenMigration generates migation file.
type CmdGenMigration struct {
	FS afero.Afero
	VersionGenerator
}

type genMigrationColumn struct {
	dbname string
	dbtype string
	null   bool
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
	err := c.Arguments().Validate(o.NamedInputs)
	if err != nil {
		return nil, err
	}
	kcls := []ToKCL{}
	kclsAfter := []ToKCL{}
	for _, v := range o.PositionalInputs {
		splits := strings.Split(v, ":")
		if len(splits) == 1 {
			switch splits[0] {
			case "id":
				kcls = append(kcls, &DDLColumn{Name: "id", Type: "bigint", Null: true, Identity: &DDLIdentity{}})
				kclsAfter = append(kclsAfter, &DDLPrimaryKey{Columns: []string{"id"}})
			case "timestamps":
				kcls = append(kcls, &DDLColumn{Name: "created_at", Type: "timestamptz", Null: false})
				kcls = append(kcls, &DDLColumn{Name: "updated_at", Type: "timestamptz", Null: false})
			default:
				return nil, fmt.Errorf("Invalid column: %s", splits[0])
			}
		} else {
			kcls = append(kcls, &DDLColumn{Name: splits[0], Type: splits[1], Null: true})
		}
	}

	kcls = append(kcls, kclsAfter...)
	return c.run(ctx, o.NamedInputs, kcls)
}

func (c *CmdGenMigration) run(ctx context.Context, inputs NamedInputs, columns []ToKCL) (ResponseGenMigration, error) {
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
	for _, col := range columns {
		sb := strings.Builder{}
		col.ToKCL(&sb)
		lines := strings.Split(sb.String(), "\n")
		for _, line := range lines {
			buf.WriteString("    ")
			buf.WriteString(line)
			buf.WriteString("\n")
		}
	}
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

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

func (c *CmdGenMigration) Do(ctx context.Context, o CmdOpts) (any, error) {
	err := c.Arguments().Validate(o.NamedInputs)
	if err != nil {
		return nil, err
	}
	return c.run(ctx, o)
}

func (c *CmdGenMigration) run(ctx context.Context, o CmdOpts) (ResponseGenMigration, error) {
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

	columns := []*DDLColumn{}
	pks := []*DDLPrimaryKey{}
	for _, v := range o.PositionalInputs {
		splits := strings.Split(v, ":")
		if len(splits) == 1 {
			switch splits[0] {
			case "id":
				columns = append(columns, &DDLColumn{Name: "id", Type: "bigint", Null: true, Identity: &DDLIdentity{}})
				pks = append(pks, &DDLPrimaryKey{Columns: []string{"id"}})
			case "timestamps":
				columns = append(columns, &DDLColumn{Name: "created_at", Type: "timestamptz", Null: false})
				columns = append(columns, &DDLColumn{Name: "updated_at", Type: "timestamptz", Null: false})
			default:
				return result, fmt.Errorf("Invalid column: %s", splits[0])
			}
		} else {
			columns = append(columns, &DDLColumn{Name: splits[0], Type: splits[1], Null: true})
		}
	}

	ref := o.NamedInputs["name"].(string)
	table := ref
	words := strings.SplitN(ref, "_", 2)
	if len(words) == 2 && words[0] == "create" {
		table = words[1]
	}

	version := c.VersionGenerator.Next()
	migration := &Migration{
		RefName: ref,
		Version: version,
		Up: MigrationUpOrDown{
			CreateTables: []*DDLCreateTable{
				{
					Name:        table,
					Columns:     columns,
					PrimaryKeys: pks,
				},
			},
		},
		Down: MigrationUpOrDown{
			DropTables: []*DDLDropTable{
				{
					Name: table,
				},
			},
		},
	}

	result.Ref = fmt.Sprint("migration.", ref)
	result.Path = filepath.Join(dir, fmt.Sprint(version, "_", ref, krabenv.Ext()))

	buf := bytes.Buffer{}
	migration.ToKCL(&buf)

	c.FS.WriteFile(result.Path, buf.Bytes(), 0644)

	return result, nil
}

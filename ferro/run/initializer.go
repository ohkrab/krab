package run

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/ohkrab/krab/ferro/config"
	"github.com/ohkrab/krab/ferro/run/generators"
)

type Initializer struct {
	fs                 *config.Filesystem
	timestampGenerator generators.TimestampGenerator
	renderer           generators.TemplateRenderer
}

func NewInitializer(fs *config.Filesystem, timestampGenerator generators.TimestampGenerator) *Initializer {
	return &Initializer{
		fs:                 fs,
		timestampGenerator: timestampGenerator,
		renderer:           generators.TemplateRenderer{},
	}
}

type InitializeOptions struct {
}

func (init *Initializer) Initialize(ctx context.Context, opts InitializeOptions) error {
	err := init.fs.MkdirAll([]string{".ferro", "migrations", "public"})
	if err != nil {
		return err
	}

	prefix := filepath.Join(".ferro", "migrations")

	now := init.timestampGenerator.Next()

	renderedSet := init.renderer.Render(
		`---
apiVersion: migrations/v1
kind: MigrationSet
metadata:
  name: public
spec:
  namespace:
    name: public
  migrations:
    - create_hello_world
    # - another_migration
`,
	)

	err = init.fs.TouchFile(filepath.Join(prefix, "public.fyml"), []byte(renderedSet))
	if err != nil {
		return err
	}

	renderedMigration := init.renderer.Render(
		`---
apiVersion: migrations/v1
kind: Migration
metadata:
  name: create_hello_world
spec:
  version: "` + now.String() + `"
  run:
    up:
      sql: |
        CREATE TABLE hello_world (
          id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
          name VARCHAR(255) NOT NULL
        );
    down:
      sql: DROP TABLE hello_world;
`,
	)

	err = init.fs.TouchFile(filepath.Join(prefix, "public", fmt.Sprintf("%s_create_hello_world.fyml", now)), []byte(renderedMigration))
	if err != nil {
		return err
	}

	return nil
}

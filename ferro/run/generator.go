package run

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/ohkrab/krab/ferro/config"
	"github.com/ohkrab/krab/ferro/run/generators"
	"github.com/ohkrab/krab/tpls"
)

type Generator struct {
	fs                 *config.Filesystem
	tpls               *tpls.Templates
	timestampGenerator generators.TimestampGenerator
}

func NewGenerator(fs *config.Filesystem, tpls *tpls.Templates, timestampGenerator generators.TimestampGenerator) *Generator {
	return &Generator{fs: fs, tpls: tpls, timestampGenerator: timestampGenerator}
}

type GenerateInitOptions struct {
}

func (g *Generator) GenInit(ctx context.Context, opts GenerateInitOptions) error {
	if g.fs.Exists(".ferro") {
		return fmt.Errorf(".ferro directory already exists")
	}

	err := g.fs.MkdirAll([]string{".ferro", "migrations", "public"})
	if err != nil {
		return err
	}

	prefix := filepath.Join(".ferro", "migrations")

	now := g.timestampGenerator.Next()

	renderedSet, err := g.tpls.RenderEmbedded("set", map[string]any{
		"Name": "public",
	})
	if err != nil {
		return err
	}

	err = g.fs.TouchFile(filepath.Join(prefix, "public.fyml"), renderedSet)
	if err != nil {
		return err
	}

	renderedMigration, err := g.tpls.RenderEmbedded("migration", map[string]any{
		"Name":    "create_hello_world",
		"Version": now.String(),
		"Up":      "CREATE TABLE hello_world (id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY, name VARCHAR(255) NOT NULL);",
		"Down":    "DROP TABLE hello_world;",
	})
	if err != nil {
		return err
	}

	err = g.fs.TouchFile(filepath.Join(prefix, "public", fmt.Sprintf("%s_create_hello_world.fyml", now)), renderedMigration)
	if err != nil {
		return err
	}

	return nil
}

type GenerateMigrationOptions struct {
}

func (g *Generator) GenMigration(ctx context.Context, opts GenerateMigrationOptions) error {
	err := g.fs.MkdirAll([]string{".ferro", "migrations", "public"})
	if err != nil {
		return err
	}

	prefix := filepath.Join(".ferro", "migrations")

	now := g.timestampGenerator.Next()

	renderedMigration, err := g.tpls.RenderEmbedded("migration", map[string]any{
		"Name":    "create_hello_world",
		"Version": now.String(),
		"Up":      "CREATE TABLE hello_world (id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY, name VARCHAR(255) NOT NULL);",
		"Down":    "DROP TABLE hello_world;",
	})
	if err != nil {
		return err
	}

	err = g.fs.TouchFile(filepath.Join(prefix, "public", fmt.Sprintf("%s_create_hello_world.fyml", now)), renderedMigration)
	if err != nil {
		return err
	}

	return nil
}

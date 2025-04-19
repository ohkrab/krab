package run

import (
	"context"
	"io"

	"github.com/ohkrab/krab/ferro/config"
	"github.com/ohkrab/krab/ferro/plugin"
)

type Migrator struct {
	fs *config.Filesystem
}

func NewMigrator(fs *config.Filesystem) *Migrator {
	return &Migrator{
		fs: fs,
	}
}

type MigrateInitOptions struct {
}

func (m *Migrator) MigrateInit(ctx context.Context, config *config.Config, opts MigrateInitOptions) error {
	return nil
}

type MigrateAuditOptions struct {
	Driver     plugin.Driver
	Output     io.Writer
	FilterLast int
}

func (m *Migrator) MigrateAudit(ctx context.Context, config *config.Config, opts MigrateAuditOptions) error {
	return nil
}

type MigrateUpOptions struct {
	Driver plugin.Driver
	Set    string
}

func (m *Migrator) MigrateUp(ctx context.Context, config *config.Config, opts MigrateUpOptions) error {
	return nil
}

type MigrateDownOptions struct {
	Driver  plugin.Driver
	Set     string
	Version string
}

func (m *Migrator) MigrateDown(ctx context.Context, config *config.Config, opts MigrateDownOptions) error {
	return nil
}

type MigrateStatusOptions struct {
	Driver plugin.Driver
	Set    string
}

func (m *Migrator) MigrateStatus(ctx context.Context, config *config.Config, opts MigrateStatusOptions) error {
	return nil
}

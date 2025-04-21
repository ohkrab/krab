package run

import (
	"context"

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

type MigrateAuditOptions struct {
	Driver      plugin.DriverInstance
	Set         string
	FilterLastN uint
}

func (m *Migrator) MigrateAudit(ctx context.Context, config *config.Config, opts MigrateAuditOptions) error {
	return nil
}

type MigrateUpOptions struct {
	Driver plugin.DriverInstance
	Set    string
}

func (m *Migrator) MigrateUp(ctx context.Context, config *config.Config, opts MigrateUpOptions) error {
	return nil
}

type MigrateDownOptions struct {
	Driver  plugin.DriverInstance
	Set     string
	Version string
}

func (m *Migrator) MigrateDown(ctx context.Context, config *config.Config, opts MigrateDownOptions) error {
	return nil
}

type MigrateStatusOptions struct {
	Driver plugin.DriverInstance
	Set    string
}

func (m *Migrator) MigrateStatus(ctx context.Context, config *config.Config, opts MigrateStatusOptions) error {
	nav := NewNavigator(opts.Driver, config)
	conn, close, err := nav.Open(ctx)
	if err != nil {
		return err
	}
	defer close()

	err = nav.Ready(ctx, conn)
	if err != nil {
		return err
	}

	err = nav.Drive(ctx, conn, func() error {
		return nil
	})
	return err
}

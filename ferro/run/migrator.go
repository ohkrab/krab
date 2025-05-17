package run

import (
	"context"
	"fmt"

	"github.com/ohkrab/krab/ferro/config"
	"github.com/ohkrab/krab/ferro/plugin"
	"github.com/ohkrab/krab/fmtx"
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
	Set         *config.MigrationSet
	FilterLastN uint
}

func (m *Migrator) MigrateAudit(ctx context.Context, config *config.Config, opts MigrateAuditOptions) (*Audited, error) {
	fmtx.WriteSuccess("Executing Migrate.Audit with Driver=%s, Set=%s", opts.Driver.Config.Metadata.Name, opts.Set.Metadata.Name)

	nav := NewNavigator(opts.Driver, config, plugin.DriverExecutionContext{
		Prefix: opts.Set.Spec.Namespace.Prefix,
		Schema: opts.Set.Spec.Namespace.Schema,
	})
	conn, close, err := nav.Open(ctx)
	if err != nil {
		return nil, err
	}
	defer close()

	err = nav.Ready(ctx, conn)
	if err != nil {
		return nil, err
	}

    var audited *Audited
	err = nav.Drive(ctx, conn, func() error {
        audited, err = nav.ComputeState(ctx, conn)
		if err != nil {
			return err
		}

		return nil
	})
    if err != nil {
        return nil, err
    }

	return audited, nil
}

type MigrateUpOptions struct {
	Driver plugin.DriverInstance
	Set    *config.MigrationSet
}

func (m *Migrator) MigrateUp(ctx context.Context, config *config.Config, opts MigrateUpOptions) error {
	fmtx.WriteSuccess("Executing Migrate.Up with Driver=%s, Set=%s", opts.Driver.Config.Metadata.Name, opts.Set.Metadata.Name)

	return fmt.Errorf("not implemented")
}

type MigrateDownOptions struct {
	Driver  plugin.DriverInstance
	Set     *config.MigrationSet
	Version string
}

func (m *Migrator) MigrateDown(ctx context.Context, config *config.Config, opts MigrateDownOptions) error {
	fmtx.WriteSuccess("Executing Migrate.Down with Driver=%s, Set=%s, Version=%s", opts.Driver.Config.Metadata.Name, opts.Set.Metadata.Name, opts.Version)

	return fmt.Errorf("not implemented")
}

type MigrateStatusOptions struct {
	Driver plugin.DriverInstance
	Set    *config.MigrationSet
}

func (m *Migrator) MigrateStatus(ctx context.Context, config *config.Config, opts MigrateStatusOptions) error {
	fmtx.WriteSuccess("Executing Migrate.Status with Driver=%s, Set=%s", opts.Driver.Config.Metadata.Name, opts.Set.Metadata.Name)

	nav := NewNavigator(opts.Driver, config, plugin.DriverExecutionContext{
		Prefix: opts.Set.Spec.Namespace.Prefix,
		Schema: opts.Set.Spec.Namespace.Schema,
	})
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
		_, err := nav.ComputeState(ctx, conn)
		if err != nil {
			return err
		}

		return nil
	})
	return err
}

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

type MigrateAuditResult struct {
}

func (m *Migrator) MigrateAudit(ctx context.Context, config *config.Config, opts MigrateAuditOptions) (*MigrateAuditResult, error) {
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

	err = nav.Drive(ctx, conn, func() error {
		_, err := nav.ComputeState(ctx, conn)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &MigrateAuditResult{}, nil
}

type MigrateUpOptions struct {
	Driver plugin.DriverInstance
	Set    *config.MigrationSet
}

type MigrateUpResult struct {
}

func (m *Migrator) MigrateUp(ctx context.Context, config *config.Config, opts MigrateUpOptions) (*MigrateUpResult, error) {
	fmtx.WriteSuccess("Executing Migrate.Up with Driver=%s, Set=%s", opts.Driver.Config.Metadata.Name, opts.Set.Metadata.Name)

	return nil, fmt.Errorf("not implemented")
}

type MigrateDownOptions struct {
	Driver  plugin.DriverInstance
	Set     *config.MigrationSet
	Version string
}

type MigrateDownResult struct {
}

func (m *Migrator) MigrateDown(ctx context.Context, config *config.Config, opts MigrateDownOptions) (*MigrateDownResult, error) {
	fmtx.WriteSuccess("Executing Migrate.Down with Driver=%s, Set=%s, Version=%s", opts.Driver.Config.Metadata.Name, opts.Set.Metadata.Name, opts.Version)

	return nil, fmt.Errorf("not implemented")
}

type MigrateStatusOptions struct {
	Driver plugin.DriverInstance
	Set    *config.MigrationSet
}

type MigrateStatusResult struct {
	Rows []MigrationStatusResultSingle
}

type MigrationStatusResultSingle struct {
	Migration string
	Status    string
	Version   string
}

func (m *Migrator) MigrateStatus(ctx context.Context, config *config.Config, opts MigrateStatusOptions) (*MigrateStatusResult, error) {
	fmtx.WriteSuccess("Executing Migrate.Status with Driver=%s, Set=%s", opts.Driver.Config.Metadata.Name, opts.Set.Metadata.Name)

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

	rows := make([]MigrationStatusResultSingle, 0)
    auditedSet := audited.EnsureMigrationSet(opts.Set.Metadata.Name)
    for _, name := range opts.Set.Spec.Migrations {
        row := MigrationStatusResultSingle{
            Migration: name,
            Status:    "unknown",
            Version:  "",
        }
        specMigration, ok := config.Migrations[name]
        if ok {
            row.Version = specMigration.Spec.Version
            row.Status = "pending"
        } else {
            row.Version = "<missing>"
            row.Status = "<missing>"
        }
        auditedMigration, ok := auditedSet.Migrations[name]
        if ok {
            row.Status = auditedMigration.Status
        }
        rows = append(rows, row)
    }

	return &MigrateStatusResult{Rows: rows}, nil
}

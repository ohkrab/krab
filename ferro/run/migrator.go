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

func (m *Migrator) MigrateAudit(ctx context.Context, cfg *config.Config, opts MigrateAuditOptions) (*MigrateAuditResult, error) {
	fmtx.WriteSuccess("Executing Migrate.Audit with Driver=%s, Set=%s", opts.Driver.Config.Metadata.Name, opts.Set.Metadata.Name)

	nav := NewNavigator(opts.Driver, cfg, plugin.DriverExecutionContext{
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

func (m *Migrator) MigrateUp(ctx context.Context, cfg *config.Config, opts MigrateUpOptions) (*MigrateUpResult, error) {
	fmtx.WriteSuccess("Executing Migrate.Up with Driver=%s, Set=%s", opts.Driver.Config.Metadata.Name, opts.Set.Metadata.Name)

	nav := NewNavigator(opts.Driver, cfg, plugin.DriverExecutionContext{
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

	var result MigrateUpResult
	err = nav.Drive(ctx, conn, func() error {
		audited, err := nav.ComputeState(ctx, conn)
		if err != nil {
			return err
		}
		auditedSet := audited.EnsureMigrationSet(opts.Set.Metadata.Name)
		pendingMigrations := make([]*config.Migration, 0)
		for _, name := range opts.Set.Spec.Migrations {
			specMigration, ok := cfg.Migrations[name]
			if ok {
				auditedMigration, ok := auditedSet.Migrations[name]
				if ok {
					if auditedMigration.Status == AuditStatusFailed {
                        return fmt.Errorf("exec: Migration %s is in a failed state, please fix the migration before proceeding", name)
					}
				} else {
					pendingMigrations = append(pendingMigrations, specMigration)
				}
			} else {
				return fmt.Errorf("exec: Migration %s not found in config", name)
			}
		}

		for _, migration := range pendingMigrations {
			fmtx.WriteInfo("Executing migration: %s", migration.Metadata.Name)
            // mark as started
            // execute the SQL migration
            // mark as completed or failed
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &result, nil
}

type MigrateDownOptions struct {
	Driver  plugin.DriverInstance
	Set     *config.MigrationSet
	Version string
}

type MigrateDownResult struct {
}

func (m *Migrator) MigrateDown(ctx context.Context, cfg *config.Config, opts MigrateDownOptions) (*MigrateDownResult, error) {
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

func (m *Migrator) MigrateStatus(ctx context.Context, cfg *config.Config, opts MigrateStatusOptions) (*MigrateStatusResult, error) {
	fmtx.WriteSuccess("Executing Migrate.Status with Driver=%s, Set=%s", opts.Driver.Config.Metadata.Name, opts.Set.Metadata.Name)

	nav := NewNavigator(opts.Driver, cfg, plugin.DriverExecutionContext{
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
			Version:   "",
		}
		specMigration, ok := cfg.Migrations[name]
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

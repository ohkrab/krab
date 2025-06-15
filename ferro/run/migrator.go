package run

import (
	"context"
	"fmt"
	"time"

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
	Logs []plugin.DriverAuditLog
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

	result := &MigrateAuditResult{
		Logs: make([]plugin.DriverAuditLog, 0),
	}

	err = nav.Drive(ctx, conn, func() error {
		audited, err := nav.ComputeState(ctx, conn)
		if err != nil {
			return err
		}

		for _, log := range audited.Raw {
			if log.Data["set"] == opts.Set.Metadata.Name {
				result.Logs = append(result.Logs, log)
			}
		}

        if opts.FilterLastN > 0 {
            result.Logs = result.Logs[len(result.Logs)-int(opts.FilterLastN):]
        }

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

type MigrateUpOptions struct {
	Driver plugin.DriverInstance
	Set    *config.MigrationSet
}

type MigrateUpResult struct {
	WasPending int
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
				auditedMigration, ok := auditedSet.Migrations[MigrationVersion(specMigration.Spec.Version)]
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

		result.WasPending = len(pendingMigrations)

		for _, pending := range pendingMigrations {
			started := plugin.DriverAuditLog{
				ID:        audited.LastID + 1,
				AppliedAt: time.Now().UTC(),
				Event:     MigrationUpStartedEvent,
				Data: map[string]any{
					"set":       opts.Set.Metadata.Name,
					"migration": pending.Metadata.Name,
					"version":   pending.Spec.Version,
				},
				Metadata: map[string]any{},
			}
			err := nav.Mark(ctx, conn, started)
			if err != nil {
				return fmt.Errorf("exec: Failed to mark migration `%s` as started: %w", pending.Metadata.Name, err)
			}
			fmtx.WriteInfo("Executing migration: %s", pending.Metadata.Name)

			err = nav.WithTx(ctx, conn, true, func(query plugin.DriverQuery) error {
				return query.Exec(ctx, pending.Spec.Run.Up.Sql)
			})

			stopped := plugin.DriverAuditLog{
				ID:        started.ID + 1,
				AppliedAt: time.Now().UTC(),
				Event:     "",
				Data: map[string]any{
					"set":       opts.Set.Metadata.Name,
					"migration": pending.Metadata.Name,
					"version":   pending.Spec.Version,
				},
				Metadata: map[string]any{},
			}
			if err != nil {
				stopped.Event = MigrationUpFailedEvent
				stopped.Metadata["error"] = err.Error()
			} else {
				stopped.Event = MigrationUpCompletedEvent
			}
			err = nav.Mark(ctx, conn, stopped)
			if err != nil {
				return fmt.Errorf("critical(inconsistent state): Failed to mark migration `%s` as completed/failed: %w", pending.Metadata.Name, err)
			}
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
		res, err := nav.ComputeState(ctx, conn)
		if err != nil {
			return err
		}
		audited = res

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
		auditedMigration, ok := auditedSet.Migrations[MigrationVersion(specMigration.Spec.Version)]
		if ok {
			row.Status = auditedMigration.Status
		}
		rows = append(rows, row)
	}

	return &MigrateStatusResult{Rows: rows}, nil
}

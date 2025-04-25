package run

import (
	"context"
	"fmt"
	"os"

	"github.com/ohkrab/krab/ferro/config"
	"github.com/ohkrab/krab/ferro/plugin"
	"github.com/ohkrab/krab/fmtx"
)

// Navigator abstracts the flow of the driver.
type Navigator struct {
	driver plugin.DriverInstance
	config *config.Config
}

type Audited struct {
	Sets map[string]*AuditedSet
}

type AuditedSet struct {
	Migrations map[string]*AuditedMigration
	Status     string
}

type AuditedMigration struct {
	Version string
	Status  string
}

func NewNavigator(driver plugin.DriverInstance, config *config.Config) *Navigator {
	return &Navigator{
		driver: driver,
		config: config,
	}
}

func (n *Navigator) Open(ctx context.Context) (plugin.DriverConnection, func(), error) {
	conn, err := n.driver.Driver.Connect(ctx, n.driver.Config.Spec.Config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	disconnect := func() {
		err := n.driver.Driver.Disconnect(ctx, conn)
		if err != nil {
			fmtx.WriteError(fmt.Sprintf("failed to disconnect from database: %s", err))
			os.Exit(2)
		}
	}
	return conn, disconnect, nil
}

func (n *Navigator) Ready(ctx context.Context, conn plugin.DriverConnection) error {
	if err := conn.UpsertAuditLogTable(ctx, plugin.DriverExecutionContext{}); err != nil {
		return fmt.Errorf("failed to upsert audit log table: %w", err)
	}
	return nil
}

func (n *Navigator) Drive(ctx context.Context, conn plugin.DriverConnection, run func() error) error {
	err := conn.LockAuditLog(ctx, plugin.DriverExecutionContext{})
	if err != nil {
		return fmt.Errorf("failed to lock audit log: %w", err)
	}
	defer func() {
		err := conn.UnlockAuditLog(ctx, plugin.DriverExecutionContext{})
		if err != nil {
			fmtx.WriteError(fmt.Sprintf("failed to unlock audit log: %s", err))
			if err := n.driver.Driver.Disconnect(ctx, conn); err != nil {
				fmtx.WriteError(fmt.Sprintf("failed to disconnect from database: %s", err))
			}
			os.Exit(2)
		}
	}()
	err = run()
	if err != nil {
		return fmt.Errorf("driver failed to run: %w", err)
	}
	return nil
}

func (n *Navigator) ComputeState(ctx context.Context, conn plugin.DriverConnection) (*Audited, error) {
	logs, err := conn.ReadAuditLogs(ctx, plugin.DriverExecutionContext{})
	if err != nil {
		return nil, fmt.Errorf("failed to read audit logs: %w", err)
	}

	audited := &Audited{
		Sets: make(map[string]*AuditedSet),
	}

	set := ""

	for _, log := range logs {
		switch log.Event {
		case "migration_sets.up.started":
			set = log.GetMetadata("name")
			_, exists := audited.Sets[set]
			if !exists {
				audited.Sets[set] = &AuditedSet{
					Migrations: make(map[string]*AuditedMigration),
					Status:     "started",
				}
			}

		case "migration_sets.up.completed":
			set = log.GetMetadata("name")
			audited.Sets[set].Status = "completed"

		case "migration_sets.up.failed":
			set = log.GetMetadata("name")
			audited.Sets[set].Status = "failed"

		case "migrations.started":
			if set == "" {
				return nil, fmt.Errorf("FATAL: audit log is broken/out of order")
			}
			version := log.GetData("version")
			audited.Sets[set].Migrations[version] = &AuditedMigration{
				Version: version,
				Status:  "started",
			}

		case "migrations.completed":
			if set == "" {
				return nil, fmt.Errorf("FATAL: audit log is broken/out of order")
			}
			version := log.GetData("version")
			audited.Sets[set].Migrations[version].Status = "completed"

		case "migrations.failed":
			if set == "" {
				return nil, fmt.Errorf("FATAL: audit log is broken/out of order")
			}
			version := log.GetData("version")
			audited.Sets[set].Migrations[version].Status = "failed"
		}
	}

	return audited, nil
}

// lock_id, api_version,kind,applied_at,event,data,metadata
// 1, migrations/v1, MigrationSet, 20250505T1200Z, migration_sets.up.started, {}, {name: "public"}
// 2, migrations/v1, Migration, 20250505T1200Z, migrations.started, {version: "202006_01"}, {name: "add_tenants"}
// 3, migrations/v1, Migration, 20250505T1200Z, migrations.completed, {version: "202006_01"}, {name: "add_tenants"}
// 4, migrations/v1, Migration, 20250505T1200Z, migrations.failed, {version: "202006_01"}, {name: "add_tenants"}
// 5, migrations/v1, MigrationSet, 20250505T1200Z, migration_sets.up.completed, {}, {name: "public"}

// 12, migrations/v1, MigrationSet, 20250505T1200Z, migration_sets.up.started, {args: {schema: "animals"}}, {name: "tenant"}
// 13, migrations/v1, Migration, 20250505T1200Z, migrations.started, {version: "v1"}, {name: "create_kinds"}
// 14, migrations/v1, Migration, 20250505T1200Z, migrations.completed, {version: "v1"}, {name: "create_kinds"}
// 15, migrations/v1, Migration, 20250505T1200Z, migrations.started, {version: "v2"}, {name: "create_countries"}
// 16, migrations/v1, Migration, 20250505T1200Z, migrations.completed, {version: "v2"}, {name: "create_countries"}
// 17, migrations/v1, MigrationSet, 20250505T1200Z, migration_sets.up.completed, {args: {schema: "animals"}, {name: "tenant"}
// 18, migrations/v1, MigrationSet, 20250505T1200Z, migration_sets.renamed, {from: "public", to: "brands"}, {}

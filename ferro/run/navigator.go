package run

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ohkrab/krab/ferro/config"
	"github.com/ohkrab/krab/ferro/plugin"
	"github.com/ohkrab/krab/fmtx"
)

// Navigator abstracts the flow of the driver.
type Navigator struct {
	driver  plugin.DriverInstance
	config  *config.Config
	execCtx plugin.DriverExecutionContext
}

type Audited struct {
	Sets   map[string]*AuditedMigrationSet
	LastID uint64
}

type AuditedMigrationSet struct {
	Migrations map[string]*AuditedMigration
}

type AuditedMigration struct {
	Version string
	Status  string
}

const (
	AuditStatusStarted   = "started"
	AuditStatusCompleted = "completed"
	AuditStatusFailed    = "failed"

	MigrationUpStartedEvent     = "migration.up.started"
	MigrationUpCompletedEvent   = "migration.up.completed"
	MigrationUpFailedEvent      = "migration.up.failed"
	MigrationDownStartedEvent   = "migration.down.started"
	MigrationDownCompletedEvent = "migration.down.completed"
	MigrationDownFailedEvent    = "migration.down.failed"
)

func NewNavigator(driver plugin.DriverInstance, config *config.Config, execCtx plugin.DriverExecutionContext) *Navigator {
	return &Navigator{
		driver:  driver,
		config:  config,
		execCtx: execCtx,
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
	if err := conn.UpsertAuditLogTable(ctx, n.execCtx); err != nil {
		return fmt.Errorf("failed to upsert audit log table: %w", err)
	}
	if err := conn.UpsertAuditLockTable(ctx, n.execCtx); err != nil {
		return fmt.Errorf("failed to upsert audit lock table: %w", err)
	}
	return nil
}

func (n *Navigator) Drive(ctx context.Context, conn plugin.DriverConnection, run func() error) error {
	lock := plugin.DriverAuditLock{
		ID:       plugin.DriverAuditLockIDForMigrations,
		LockedAt: time.Now().UTC(),
		LockedBy: "cli",
		Data:     make(map[string]any),
	}
	err := conn.LockAuditLog(ctx, n.execCtx, lock)
	if err != nil {
		return fmt.Errorf("failed to lock audit log: %w", err)
	}
	defer func() {
		err := conn.UnlockAuditLog(ctx, n.execCtx, lock)
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

func (n *Navigator) Mark(ctx context.Context, conn plugin.DriverConnection, log plugin.DriverAuditLog) error {
    if log.Event == "" {
        return fmt.Errorf("fatal: audit log .Event cannot be empty (this is a bug that should be reported)")
    }
	if err := conn.AppendAuditLog(ctx, n.execCtx, log); err != nil {
		return fmt.Errorf("fatal: failed to append audit log: %w", err)
	}
	return nil
}

func (n *Navigator) WithTx(ctx context.Context, conn plugin.DriverConnection, withTx bool, yield func(query plugin.DriverQuery) error) error {
	tx := conn.Query(n.execCtx)

	if withTx {
		newTx, err := tx.Begin(ctx)
        tx = newTx
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}
	}

	err := yield(tx)

	if withTx {
		if err == nil {
			if err := tx.Commit(ctx); err != nil {
				return fmt.Errorf("failed to commit transaction: %w", err)
			}
		} else {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				return fmt.Errorf("failed to rollback transaction: %w", rollbackErr)
			}
		}
	}

    if err != nil {
        return fmt.Errorf("exec: %w", err)
    }

	return nil
}

func (n *Navigator) ComputeState(ctx context.Context, conn plugin.DriverConnection) (*Audited, error) {
	logs, err := conn.ReadAuditLogs(ctx, n.execCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to read audit logs: %w", err)
	}

	audited := &Audited{
		Sets: make(map[string]*AuditedMigrationSet),
	}

	for _, log := range logs {
		switch log.Event {
		case MigrationUpStartedEvent:
			set := audited.EnsureMigrationSet(log.GetData("set"))
			migration := set.EnsureMigration(log.GetData("version"))
			migration.Status = AuditStatusStarted

		case MigrationUpCompletedEvent:
			set := audited.EnsureMigrationSet(log.GetData("set"))
			migration := set.EnsureMigration(log.GetData("version"))
			migration.Status = AuditStatusCompleted

		case MigrationUpFailedEvent:
			set := audited.EnsureMigrationSet(log.GetData("set"))
			migration := set.EnsureMigration(log.GetData("version"))
			migration.Status = AuditStatusFailed

		case MigrationDownStartedEvent:
			set := audited.EnsureMigrationSet(log.GetData("set"))
			migration := set.EnsureMigration(log.GetData("version"))
			migration.Status = AuditStatusStarted

		case MigrationDownCompletedEvent:
			set := audited.EnsureMigrationSet(log.GetData("set"))
			set.DeleteMigration(log.GetData("version"))

		case MigrationDownFailedEvent:
			continue
		}
	}

	return audited, nil
}

func (a *Audited) EnsureMigrationSet(name string) *AuditedMigrationSet {
	set, exists := a.Sets[name]
	if !exists {
		set = &AuditedMigrationSet{
			Migrations: make(map[string]*AuditedMigration),
		}
		a.Sets[name] = set
	}
	return set
}

func (a *AuditedMigrationSet) DeleteMigration(version string) {
	delete(a.Migrations, version)
}

func (a *AuditedMigrationSet) EnsureMigration(version string) *AuditedMigration {
	migration, exists := a.Migrations[version]
	if !exists {
		migration = &AuditedMigration{
			Version: version,
		}
		a.Migrations[version] = migration
	}
	return migration
}

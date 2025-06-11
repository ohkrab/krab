package plugin

import (
	"context"
	"fmt"
	"time"

	"github.com/ohkrab/krab/ferro/config"
)

const (
	DriverAuditLogTableName        = "_ferro_audit_log"
	DriverAuditLockTableName       = "_ferro_audit_lock"
	DriverAuditLockIDForMigrations = "migrataion"
)

const (
	DriverAuditColumnString int = iota + 1
	DriverAuditColumnInt64
	DriverAuditColumnTime
	DriverAuditColumnJSON
)

var (
	// audit log columns
	DriverAuditColumnID        = DriverAuditColumn{PrimaryKey: true, Name: "id", Type: DriverAuditColumnInt64, Nullable: false}
	DriverAuditColumnAppliedAt = DriverAuditColumn{PrimaryKey: false, Name: "applied_at", Type: DriverAuditColumnTime, Nullable: false}
	DriverAuditColumnEvent     = DriverAuditColumn{PrimaryKey: false, Name: "event", Type: DriverAuditColumnString, Nullable: false}
	DriverAuditColumnData      = DriverAuditColumn{PrimaryKey: false, Name: "data", Type: DriverAuditColumnJSON, Nullable: false}
	DriverAuditColumnMetadata  = DriverAuditColumn{PrimaryKey: false, Name: "metadata", Type: DriverAuditColumnJSON, Nullable: true}

	// lock columns
	DriverAuditLockColumnID       = DriverAuditColumn{PrimaryKey: true, Name: "id", Type: DriverAuditColumnString, Nullable: false}
	DriverAuditLockColumnLockedAt = DriverAuditColumn{PrimaryKey: false, Name: "locked_at", Type: DriverAuditColumnTime, Nullable: false}
	DriverAuditLockColumnLockedBy = DriverAuditColumn{PrimaryKey: false, Name: "locked_by", Type: DriverAuditColumnString, Nullable: false}
	DriverAuditLockColumnData     = DriverAuditColumn{PrimaryKey: false, Name: "data", Type: DriverAuditColumnJSON, Nullable: false}

	ErrAuditAlreadyLocked = fmt.Errorf("audit log is already locked")
)

type DriverExecutionContext struct {
	Schema string
	Prefix string
}

type DriverAuditLog struct {
	ID        int64
	AppliedAt time.Time
	Event     string
	Data      map[string]any
	Metadata  map[string]any
}

type DriverAuditLock struct {
	ID       string
	LockedAt time.Time
	LockedBy string
	Data     map[string]any
}

func (d *DriverAuditLog) GetMetadata(key string) string {
	if d.Metadata == nil {
		return ""
	}
	if val, ok := d.Metadata[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func (d *DriverAuditLog) GetData(key string) string {
	if d.Data == nil {
		return ""
	}
	if val, ok := d.Data[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

type DriverAuditColumn struct {
	PrimaryKey bool
	Name       string
	Type       int
	Nullable   bool
}

type DriverInstance struct {
	Driver Driver
	Config *config.Driver
}

type Driver interface {
	Connect(ctx context.Context, config config.DriverConfig) (DriverConnection, error)
	Disconnect(ctx context.Context, conn DriverConnection) error
}

type DriverConnection interface {
	UpsertAuditLogTable(ctx context.Context, execCtx DriverExecutionContext) error
	UpsertAuditLockTable(ctx context.Context, execCtx DriverExecutionContext) error
	LockAuditLog(ctx context.Context, execCtx DriverExecutionContext, lock DriverAuditLock) error
	UnlockAuditLog(ctx context.Context, execCtx DriverExecutionContext, lock DriverAuditLock) error
	AppendAuditLog(ctx context.Context, execCtx DriverExecutionContext, log DriverAuditLog) error
	ReadAuditLogs(ctx context.Context, execCtx DriverExecutionContext) ([]DriverAuditLog, error)
	Query(execCtx DriverExecutionContext) DriverQuery
}

type DriverQuery interface {
	Begin(ctx context.Context) (DriverQuery, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	Exec(ctx context.Context, query string, args ...any) error
}

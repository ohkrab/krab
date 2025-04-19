package plugin

import (
	"context"
	"time"
)

const (
	DriverAuditLogTableName = ".ferro_audit_log"
)

const (
	DriverAuditLogColumnString int = iota + 1
	DriverAuditLogColumnBigInt
	DriverAuditLogColumnTime
	DriverAuditLogColumnJSON
)

var DriverAuditLogColumns = []DriverAuditLogColumn{
	{
		Unique: true,
		Name:   "id",
		Type:   DriverAuditLogColumnBigInt,
	},
	{
		Unique: false,
		Name:   "api_version",
		Type:   DriverAuditLogColumnString,
	},
	{
		Unique: false,
		Name:   "kind",
		Type:   DriverAuditLogColumnString,
	},
	{
		Unique: false,
		Name:   "applied_at",
		Type:   DriverAuditLogColumnTime,
	},
	{
		Unique: false,
		Name:   "event",
		Type:   DriverAuditLogColumnString,
	},
	{
		Unique: false,
		Name:   "data",
		Type:   DriverAuditLogColumnJSON,
	},
	{
		Unique: false,
		Name:   "metadata",
		Type:   DriverAuditLogColumnJSON,
	},
}

type DriverExecutionContext struct {
	Schema string
	Prefix string
}

type DriverAuditLog struct {
	ID         int64
	APIVersion string
	Kind       string
	AppliedAt  time.Time
	Event      string
	Data       map[string]any
	Metadata   map[string]any
}

type DriverAuditLogColumn struct {
	Unique bool
	Name   string
	Type   int
}

type Driver interface {
	Connect(ctx context.Context) (DriverConnection, error)
	Disconnect(ctx context.Context, conn DriverConnection) error
}

type DriverConnection interface {
	LockAuditLog(ctx context.Context, execCtx DriverExecutionContext) error
	UpsertAuditLogTable(ctx context.Context, execCtx DriverExecutionContext) error
	AppendAuditLog(ctx context.Context, execCtx DriverExecutionContext, log DriverAuditLog) error
	ReadAuditLogs(ctx context.Context, execCtx DriverExecutionContext) ([]DriverAuditLog, error)
	UnlockAuditLog(ctx context.Context, execCtx DriverExecutionContext) error
}

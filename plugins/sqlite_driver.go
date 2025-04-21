package plugins

import (
	"context"

	"github.com/ohkrab/krab/ferro/plugin"
)

type SQLiteDriver struct {
	plugin.Driver
}

type SQLiteDriverConnection struct {
	plugin.DriverConnection
}

func NewSQLiteDriver() *SQLiteDriver {
	return &SQLiteDriver{}
}

func (d *SQLiteDriver) Connect(ctx context.Context) (plugin.DriverConnection, error) {
	return &SQLiteDriverConnection{}, nil
}

func (d *SQLiteDriver) Disconnect(ctx context.Context, conn plugin.DriverConnection) error {
	return nil
}

func (c *SQLiteDriverConnection) LockAuditLog(ctx context.Context, execCtx plugin.DriverExecutionContext) error {
	return nil
}

func (c *SQLiteDriverConnection) UpsertAuditLogTable(ctx context.Context, execCtx plugin.DriverExecutionContext) error {
	return nil
}

func (c *SQLiteDriverConnection) AppendAuditLog(ctx context.Context, execCtx plugin.DriverExecutionContext, log plugin.DriverAuditLog) error {
	return nil
}

func (c *SQLiteDriverConnection) ReadAuditLogs(ctx context.Context, execCtx plugin.DriverExecutionContext) ([]plugin.DriverAuditLog, error) {
	return nil, nil
}

func (c *SQLiteDriverConnection) UnlockAuditLog(ctx context.Context, execCtx plugin.DriverExecutionContext) error {
	return nil
}

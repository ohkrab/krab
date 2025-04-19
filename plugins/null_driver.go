package plugins

import (
	"context"

	"github.com/ohkrab/krab/ferro/plugin"
)

type NullDriver struct {
	plugin.Driver
}

func NewNullDriver() *NullDriver {
	return &NullDriver{}
}

func (d *NullDriver) Connect(ctx context.Context) (plugin.DriverConnection, error) {
	return &NullDriverConnection{}, nil
}

func (d *NullDriver) Disconnect(ctx context.Context, conn plugin.DriverConnection) error {
	return nil
}

type NullDriverConnection struct {
	plugin.DriverConnection
}

func (c *NullDriverConnection) LockAuditLog(ctx context.Context, execCtx plugin.DriverExecutionContext) error {
	return nil
}

func (c *NullDriverConnection) UpsertAuditLogTable(ctx context.Context, execCtx plugin.DriverExecutionContext) error {
	return nil
}

func (c *NullDriverConnection) AppendAuditLog(ctx context.Context, execCtx plugin.DriverExecutionContext, log plugin.DriverAuditLog) error {
	return nil
}

func (c *NullDriverConnection) ReadAuditLogs(ctx context.Context, execCtx plugin.DriverExecutionContext) ([]plugin.DriverAuditLog, error) {
	return []plugin.DriverAuditLog{}, nil
}

func (c *NullDriverConnection) UnlockAuditLog(ctx context.Context, execCtx plugin.DriverExecutionContext) error {
	return nil
}

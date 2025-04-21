package plugins

import (
	"context"
	"fmt"

	"github.com/ohkrab/krab/ferro/config"
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

func (d *SQLiteDriver) Connect(ctx context.Context, config config.DriverConfig) (plugin.DriverConnection, error) {
	return &SQLiteDriverConnection{}, fmt.Errorf("not implemented")
}

func (d *SQLiteDriver) Disconnect(ctx context.Context, conn plugin.DriverConnection) error {
	return fmt.Errorf("not implemented")
}

func (c *SQLiteDriverConnection) LockAuditLog(ctx context.Context, execCtx plugin.DriverExecutionContext) error {
	return fmt.Errorf("not implemented")
}

func (c *SQLiteDriverConnection) UpsertAuditLogTable(ctx context.Context, execCtx plugin.DriverExecutionContext) error {
	return fmt.Errorf("not implemented")
}

func (c *SQLiteDriverConnection) AppendAuditLog(ctx context.Context, execCtx plugin.DriverExecutionContext, log plugin.DriverAuditLog) error {
	return fmt.Errorf("not implemented")
}

func (c *SQLiteDriverConnection) ReadAuditLogs(ctx context.Context, execCtx plugin.DriverExecutionContext) ([]plugin.DriverAuditLog, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *SQLiteDriverConnection) UnlockAuditLog(ctx context.Context, execCtx plugin.DriverExecutionContext) error {
	return fmt.Errorf("not implemented")
}

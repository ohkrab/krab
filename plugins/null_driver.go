package plugins

import (
	"context"
	"errors"

	"github.com/ohkrab/krab/ferro/config"
	"github.com/ohkrab/krab/ferro/plugin"
)

var ErrDriverNotSelected = errors.New("driver not selected")

type NullDriver struct {
	plugin.Driver
}

func NewNullDriver() *NullDriver {
	return &NullDriver{}
}

func (d *NullDriver) Connect(ctx context.Context, config config.DriverConfig) (plugin.DriverConnection, error) {
	return &NullDriverConnection{}, ErrDriverNotSelected
}

func (d *NullDriver) Disconnect(ctx context.Context, conn plugin.DriverConnection) error {
	return ErrDriverNotSelected
}

type NullDriverConnection struct {
	plugin.DriverConnection
}

func (c *NullDriverConnection) LockAuditLog(ctx context.Context, execCtx plugin.DriverExecutionContext) error {
	return ErrDriverNotSelected
}

func (c *NullDriverConnection) UpsertAuditLogTable(ctx context.Context, execCtx plugin.DriverExecutionContext) error {
	return ErrDriverNotSelected
}

func (c *NullDriverConnection) AppendAuditLog(ctx context.Context, execCtx plugin.DriverExecutionContext, log plugin.DriverAuditLog) error {
	return ErrDriverNotSelected
}

func (c *NullDriverConnection) ReadAuditLogs(ctx context.Context, execCtx plugin.DriverExecutionContext) ([]plugin.DriverAuditLog, error) {
	return []plugin.DriverAuditLog{}, ErrDriverNotSelected
}

func (c *NullDriverConnection) UnlockAuditLog(ctx context.Context, execCtx plugin.DriverExecutionContext) error {
	return ErrDriverNotSelected
}

package run

import (
	"context"
	"fmt"
	"os"

	"github.com/ohkrab/krab/ferro/plugin"
	"github.com/ohkrab/krab/fmtx"
)

// Navigator abstracts the flow of the driver.
type Navigator struct {
	Driver plugin.Driver
}

func NewNavigator(driver plugin.Driver) *Navigator {
	return &Navigator{
		Driver: driver,
	}
}

func (n *Navigator) Open(ctx context.Context) (plugin.DriverConnection, func(), error) {
	conn, err := n.Driver.Connect(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	disconnect := func() {
		err := n.Driver.Disconnect(ctx, conn)
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
			if err := n.Driver.Disconnect(ctx, conn); err != nil {
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

package plugins

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/ohkrab/krab/ferro/config"
	"github.com/ohkrab/krab/ferro/plugin"
)

type PostgreSQLDriver struct {
}

type PostgreSQLDriverConnection struct {
	Conn *pgx.Conn
}

func NewPostgreSQLDriver() *PostgreSQLDriver {
	return &PostgreSQLDriver{}
}

func (d *PostgreSQLDriver) Connect(ctx context.Context, config config.DriverConfig) (plugin.DriverConnection, error) {
	dsn := config.String("dsn")

	if dsn == "" {
		return nil, fmt.Errorf("config.dsn is required")
	}

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL database: %w", err)
	}
	return &PostgreSQLDriverConnection{
		Conn: conn,
	}, nil
}

func (d *PostgreSQLDriver) Disconnect(ctx context.Context, conn plugin.DriverConnection) error {
	driverConn := conn.(*PostgreSQLDriverConnection)
	err := driverConn.Conn.Close(ctx)
	if err != nil {
		return fmt.Errorf("failed to disconnect from PostgreSQL database: %w", err)
	}
	return nil
}

func (c *PostgreSQLDriverConnection) UpsertAuditLogTable(ctx context.Context, execCtx plugin.DriverExecutionContext) error {
	columns := []string{
		c.sqlColumnDefinition(&plugin.DriverAuditColumnID),
		c.sqlColumnDefinition(&plugin.DriverAuditColumnAppliedAt),
		c.sqlColumnDefinition(&plugin.DriverAuditColumnEvent),
		c.sqlColumnDefinition(&plugin.DriverAuditColumnData),
		c.sqlColumnDefinition(&plugin.DriverAuditColumnMetadata),
	}

	_, err := c.Conn.Exec(ctx,
		fmt.Sprintf(
			`CREATE TABLE IF NOT EXISTS %s (%s)`,
			c.sqlTableName(execCtx, plugin.DriverAuditLogTableName),
			strings.Join(columns, ","),
		),
	)
	if err != nil {
		return err
	}
	return nil
}

func (c *PostgreSQLDriverConnection) UpsertAuditLockTable(ctx context.Context, execCtx plugin.DriverExecutionContext) error {
	columns := []string{
		c.sqlColumnDefinition(&plugin.DriverAuditLockColumnID),
		c.sqlColumnDefinition(&plugin.DriverAuditLockColumnLockedAt),
		c.sqlColumnDefinition(&plugin.DriverAuditLockColumnLockedBy),
		c.sqlColumnDefinition(&plugin.DriverAuditLockColumnData),
	}
	_, err := c.Conn.Exec(ctx,
		fmt.Sprintf(
			`CREATE TABLE IF NOT EXISTS %s (%s)`,
			c.sqlTableName(execCtx, plugin.DriverAuditLockTableName),
			strings.Join(columns, ","),
		),
	)
	if err != nil {
		return err
	}
	return nil
}

func (c *PostgreSQLDriverConnection) AppendAuditLog(ctx context.Context, execCtx plugin.DriverExecutionContext, log plugin.DriverAuditLog) error {
	columns := []string{
		c.quoteIdentifier(plugin.DriverAuditColumnID.Name),
		c.quoteIdentifier(plugin.DriverAuditColumnAppliedAt.Name),
		c.quoteIdentifier(plugin.DriverAuditColumnEvent.Name),
		c.quoteIdentifier(plugin.DriverAuditColumnData.Name),
		c.quoteIdentifier(plugin.DriverAuditColumnMetadata.Name),
	}
	sql := `INSERT INTO %s (%s) VALUES (@id, @applied_at, @event, @data, @metadata)`
	values := pgx.NamedArgs{
		"id":         log.ID,
		"applied_at": log.AppliedAt,
		"event":      log.Event,
		"data":       log.Data,
		"metadata":   log.Metadata,
	}
	_, err := c.Conn.Exec(ctx,
		fmt.Sprintf(sql, c.sqlTableName(execCtx, plugin.DriverAuditLogTableName), strings.Join(columns, ",")),
		values,
	)
	if err != nil {
		return err
	}
	return nil
}

func (c *PostgreSQLDriverConnection) ReadAuditLogs(ctx context.Context, execCtx plugin.DriverExecutionContext) ([]plugin.DriverAuditLog, error) {
	logs := make([]plugin.DriverAuditLog, 0)
	sql := `SELECT %s FROM %s ORDER BY 1`
	columns := []string{
		c.quoteIdentifier(plugin.DriverAuditColumnID.Name),
		c.quoteIdentifier(plugin.DriverAuditColumnAppliedAt.Name),
		c.quoteIdentifier(plugin.DriverAuditColumnEvent.Name),
		c.quoteIdentifier(plugin.DriverAuditColumnData.Name),
		c.quoteIdentifier(plugin.DriverAuditColumnMetadata.Name),
	}
	rows, err := c.Conn.Query(ctx, fmt.Sprintf(sql, strings.Join(columns, ","), c.sqlTableName(execCtx, plugin.DriverAuditLogTableName)))
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var log plugin.DriverAuditLog
		err := rows.Scan(
			&log.ID,
			&log.AppliedAt,
			&log.Event,
			&log.Data,
			&log.Metadata,
		)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func (c *PostgreSQLDriverConnection) LockAuditLog(ctx context.Context, execCtx plugin.DriverExecutionContext, lock plugin.DriverAuditLock) error {
	values := pgx.NamedArgs{
		"id":        lock.ID,
		"locked_at": lock.LockedAt,
		"locked_by": lock.LockedBy,
		"data":      lock.Data,
	}
	sql := fmt.Sprintf(
		`INSERT INTO %s (%s, %s, %s, %s) VALUES (@id, @locked_at, @locked_by, @data)`,
		c.sqlTableName(execCtx, plugin.DriverAuditLockTableName),
		c.quoteIdentifier(plugin.DriverAuditLockColumnID.Name),
		c.quoteIdentifier(plugin.DriverAuditLockColumnLockedAt.Name),
		c.quoteIdentifier(plugin.DriverAuditLockColumnLockedBy.Name),
		c.quoteIdentifier(plugin.DriverAuditLockColumnData.Name),
	)
	fmt.Println(sql)

	_, err := c.Conn.Exec(ctx, sql, values)

	if err != nil {
		return err
	}

	return nil
}

func (c *PostgreSQLDriverConnection) UnlockAuditLog(ctx context.Context, execCtx plugin.DriverExecutionContext, lock plugin.DriverAuditLock) error {
	_, err := c.Conn.Exec(ctx,
		fmt.Sprintf(
			`DELETE FROM %s WHERE %s = $1`,
			c.sqlTableName(execCtx, plugin.DriverAuditLockTableName),
			c.quoteIdentifier(plugin.DriverAuditLockColumnID.Name),
		),
		lock.ID,
	)

	if err != nil {
		return err
	}
	return nil
}

func (c *PostgreSQLDriverConnection) quoteIdentifier(ident string) string {
	return pgx.Identifier{ident}.Sanitize()
}

func (c *PostgreSQLDriverConnection) sqlColumnDefinition(column *plugin.DriverAuditColumn) string {
	colDef := "UNKNOWN_TYPE"
	switch column.Type {
	case plugin.DriverAuditColumnTime:
		colDef = "TIMESTAMPTZ"
	case plugin.DriverAuditColumnString:
		colDef = "VARCHAR"
	case plugin.DriverAuditColumnInt64:
		colDef = "BIGINT"
	case plugin.DriverAuditColumnJSON:
		colDef = "JSONB"
	}
	if column.PrimaryKey {
		colDef = fmt.Sprintf("%s PRIMARY KEY", colDef)
	} else if !column.Nullable {
		colDef = fmt.Sprintf("%s NOT NULL", colDef)
	}

	return fmt.Sprintf("%s %s", pgx.Identifier{column.Name}.Sanitize(), colDef)
}

func (c *PostgreSQLDriverConnection) sqlTableName(execCtx plugin.DriverExecutionContext, table string) string {
	fullTableName := pgx.Identifier{execCtx.Prefix + table}
	if execCtx.Schema != "" {
		fullTableName = pgx.Identifier{execCtx.Schema, fullTableName[0]}
	}
	return fullTableName.Sanitize()
}

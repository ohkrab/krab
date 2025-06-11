package testcontainers

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ohkrab/krab/ferro/config"
	"github.com/ohkrab/krab/ferro/plugin"
)

type TestContainerPostgreSQLDriver struct {
}

type TestContainerPostgreSQLDriverConnection struct {
	Container *Container
	Conn      *pgx.Conn
	Close     func(ctx context.Context)
}

func NewTestContainerPostgreSQLDriver() *TestContainerPostgreSQLDriver {
	return &TestContainerPostgreSQLDriver{}
}

func (d *TestContainerPostgreSQLDriver) Connect(ctx context.Context, config config.DriverConfig) (plugin.DriverConnection, error) {
	version := config.String("version")
	user := config.String("user")
	password := config.String("password")
	db := config.String("db")
	port := config.Int("port")

	if version == "" {
		return nil, fmt.Errorf("config.version is required")
	}
	if user == "" {
		return nil, fmt.Errorf("config.user is required")
	}
	if password == "" {
		return nil, fmt.Errorf("config.password is required")
	}
	if db == "" {
		return nil, fmt.Errorf("config.db is required")
	}
	if port == 0 {
		return nil, fmt.Errorf("config.port is required")
	}

	image := fmt.Sprintf("postgres:%s-bookworm", version)
	container := &Container{
		Image: image,
		Port:  strconv.Itoa(port),
		Env: map[string]string{
			"POSTGRES_USER":     user,
			"POSTGRES_PASSWORD": password,
			"POSTGRES_DB":       db,
		},
	}
	endpoint, stop, err := container.Start(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start testcontainer: %w", err)
	}
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, endpoint, db)
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		defer stop(ctx)
		return nil, fmt.Errorf("failed to connect to testcontainer: %w", err)
	}
	return &TestContainerPostgreSQLDriverConnection{
		Container: container,
		Conn:      conn,
		Close:     stop,
	}, nil
}

func (d *TestContainerPostgreSQLDriver) Disconnect(ctx context.Context, conn plugin.DriverConnection) error {
	driverConn := conn.(*TestContainerPostgreSQLDriverConnection)
	defer driverConn.Close(ctx)

	err := driverConn.Conn.Close(ctx)
	if err != nil {
		return fmt.Errorf("failed to disconnect from testcontainer: %w", err)
	}
	return nil
}

func (c *TestContainerPostgreSQLDriverConnection) UpsertAuditLogTable(ctx context.Context, execCtx plugin.DriverExecutionContext) error {
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

func (c *TestContainerPostgreSQLDriverConnection) UpsertAuditLockTable(ctx context.Context, execCtx plugin.DriverExecutionContext) error {
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

func (c *TestContainerPostgreSQLDriverConnection) AppendAuditLog(ctx context.Context, execCtx plugin.DriverExecutionContext, log plugin.DriverAuditLog) error {
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

func (c *TestContainerPostgreSQLDriverConnection) ReadAuditLogs(ctx context.Context, execCtx plugin.DriverExecutionContext) ([]plugin.DriverAuditLog, error) {
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

func (c *TestContainerPostgreSQLDriverConnection) LockAuditLog(ctx context.Context, execCtx plugin.DriverExecutionContext, lock plugin.DriverAuditLock) error {
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

	_, err := c.Conn.Exec(ctx, sql, values)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" { // unique violation
			return plugin.ErrAuditAlreadyLocked
		}
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func (c *TestContainerPostgreSQLDriverConnection) UnlockAuditLog(ctx context.Context, execCtx plugin.DriverExecutionContext, lock plugin.DriverAuditLock) error {
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

func (c *TestContainerPostgreSQLDriverConnection) Query(execCtx plugin.DriverExecutionContext) plugin.DriverQuery {
	return &TestContainerPostgreSQLQuery{conn: c, tx: nil}
}

func (c *TestContainerPostgreSQLDriverConnection) quoteIdentifier(ident string) string {
	return pgx.Identifier{ident}.Sanitize()
}

func (c *TestContainerPostgreSQLDriverConnection) sqlColumnDefinition(column *plugin.DriverAuditColumn) string {
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

func (c *TestContainerPostgreSQLDriverConnection) sqlTableName(execCtx plugin.DriverExecutionContext, table string) string {
	fullTableName := pgx.Identifier{execCtx.Prefix + table}
	if execCtx.Schema != "" {
		fullTableName = pgx.Identifier{execCtx.Schema, fullTableName[0]}
	}
	return fullTableName.Sanitize()
}

type TestContainerPostgreSQLQuery struct {
	conn *TestContainerPostgreSQLDriverConnection
	tx   pgx.Tx
}

func (q *TestContainerPostgreSQLQuery) Exec(ctx context.Context, query string, args ...any) error {
	if q.tx != nil {
		_, err := q.tx.Exec(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	} else {
		_, err := q.conn.Conn.Exec(ctx, query, args...)
		if err != nil {
		}
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}

func (q *TestContainerPostgreSQLQuery) Begin(ctx context.Context) (plugin.DriverQuery, error) {
	tx, err := q.conn.Conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return &TestContainerPostgreSQLQuery{
		conn: q.conn,
		tx:   tx,
	}, nil
}

func (q *TestContainerPostgreSQLQuery) Commit(ctx context.Context) error {
	if q.tx == nil {
		return errors.New("no transaction to commit")
	}
	err := q.tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	q.tx = nil // reset tx after Commit
	return nil
}

func (q *TestContainerPostgreSQLQuery) Rollback(ctx context.Context) error {
	if q.tx == nil {
		return errors.New("no transaction to rollback")
	}
	err := q.tx.Rollback(ctx)
	if err != nil {
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}
	q.tx = nil // reset tx after Rollback
	return nil
}

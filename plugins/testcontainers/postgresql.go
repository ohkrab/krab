package testcontainers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/ohkrab/krab/ferro/config"
	"github.com/ohkrab/krab/ferro/plugin"
)

type TestContainerPostgreSQLDriver struct {
	plugin.Driver
}

type TestContainerPostgreSQLDriverConnection struct {
	plugin.DriverConnection

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

func (c *TestContainerPostgreSQLDriverConnection) LockAuditLog(ctx context.Context, execCtx plugin.DriverExecutionContext) error {
	// fullTableName := pgx.Identifier{execCtx.Prefix + plugin.DriverAuditLogTableName}
	// if execCtx.Schema != "" {
	// 	fullTableName = pgx.Identifier{execCtx.Schema, fullTableName[0]}
	// }
	// quotedTableName := fullTableName.Sanitize()
	// _, err := c.Conn.Exec(ctx, fmt.Sprintf("LOCK TABLE %s IN ACCESS EXCLUSIVE MODE", quotedTableName))
	// return err
	return fmt.Errorf("not lock")
}

func (c *TestContainerPostgreSQLDriverConnection) UpsertAuditLogTable(ctx context.Context, execCtx plugin.DriverExecutionContext) error {
	return fmt.Errorf("not implemented")
}

func (c *TestContainerPostgreSQLDriverConnection) AppendAuditLog(ctx context.Context, execCtx plugin.DriverExecutionContext, log plugin.DriverAuditLog) error {
	return fmt.Errorf("not implemented")
}

func (c *TestContainerPostgreSQLDriverConnection) ReadAuditLogs(ctx context.Context, execCtx plugin.DriverExecutionContext) ([]plugin.DriverAuditLog, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *TestContainerPostgreSQLDriverConnection) UnlockAuditLog(ctx context.Context, execCtx plugin.DriverExecutionContext) error {
	return fmt.Errorf("not implemented")
}

package testcontainers

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/ohkrab/krab/fmtx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type Container struct {
	Image     string
	Port      string
	Env       map[string]string
	container testcontainers.Container
	logger    *fmtx.Logger
}

func (c *Container) Start(ctx context.Context) (string, func(ctx context.Context), error) {
	req := testcontainers.ContainerRequest{
		Image:      c.Image,
		WaitingFor: wait.ForListeningPort(nat.Port(fmt.Sprintf("%s/tcp", c.Port))),
		Env:        c.Env,
		HostConfigModifier: func(hc *container.HostConfig) {
			hc.AutoRemove = true
		},
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return "", nil, err
	}
	c.container = container

	stop := func(ctx context.Context) {
		err := c.Stop(ctx)
		if err != nil {
			c.logger.WriteError("failed to stop testcontainer: %s", err)
		}
	}

	endpoint, err := c.container.Endpoint(ctx, "")
	if err != nil {
		defer stop(ctx)
		return "", nil, err
	}

	return endpoint, stop, nil
}

func (c *Container) Stop(ctx context.Context) error {
	return testcontainers.TerminateContainer(c.container)
}

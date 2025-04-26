package plugins

import (
	"fmt"

	"github.com/ohkrab/krab/ferro/plugin"
	"github.com/ohkrab/krab/plugins/testcontainers"
)

type Registry struct {
	drivers map[string]plugin.Driver
}

func New() *Registry {
	return &Registry{
		drivers: make(map[string]plugin.Driver),
	}
}

func (r *Registry) Register(name string, driver plugin.Driver) {
	if _, ok := r.drivers[name]; ok {
		panic("driver already registered: " + name)
	}
	r.drivers[name] = driver
}

func (r *Registry) Get(name string) (plugin.Driver, error) {
	driver, ok := r.drivers[name]
	if !ok {
		return nil, fmt.Errorf("driver not found: %s", name)
	}
	return driver, nil
}

func (r *Registry) RegisterAll() {
	r.Register("null", NewNullDriver())
	r.Register("sqlite", NewSQLiteDriver())
	r.Register("testcontainer/postgresql", testcontainers.NewTestContainerPostgreSQLDriver())
	r.Register("postgresql", NewPostgreSQLDriver())
}

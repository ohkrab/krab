package krab

import (
	"github.com/ohkrab/krab/krabdb"
	"github.com/ohkrab/krab/krabenv"
	"github.com/spf13/afero"
)

// CmdRegistry is a list of registred commands.
type CmdRegistry struct {
	Commands []Cmd

	FS afero.Afero
	VersionGenerator
}

// Register appends new command to registry.
func (r *CmdRegistry) Register(c Cmd) {
	r.Commands = append(r.Commands, c)
}

// RegisterAll registers all commands in the registry.
func (r *CmdRegistry) RegisterAll(config *Config, conn krabdb.Connection) {
	r.Register(&CmdVersion{})

	r.Register(&CmdGenMigration{FS: r.FS, VersionGenerator: r.VersionGenerator})

	for _, action := range config.Actions {
		action := action

		r.Register(&CmdAction{
			Action:     action,
			Connection: conn,
		})
	}

	for _, set := range config.MigrationSets {
		set := set

		r.Register(&CmdMigrateStatus{
			Set:        set,
			Connection: conn,
		})
		r.Register(&CmdMigrateDown{
			Set:        set,
			Connection: conn,
		})
		r.Register(&CmdMigrateUp{
			Set:        set,
			Connection: conn,
		})
	}

	if krabenv.Test() {
		r.Register(&CmdTestRun{
			Suite:      config.TestSuite,
			Connection: conn,
		})
	}
}

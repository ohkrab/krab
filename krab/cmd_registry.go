package krab

import (
	"github.com/ohkrab/krab/krabdb"
)

// CmdRegistry is a list of registred commands.
type CmdRegistry struct {
	Commands []Cmd
}

// Register appends new command to registry.
func (r *CmdRegistry) Register(c Cmd) {
	r.Commands = append(r.Commands, c)
}

// RegisterAll registers all commands in the registry.
func (r *CmdRegistry) RegisterAll(config *Config, conn krabdb.Connection) {
	r.Register(&CmdVersion{})

	for _, set := range config.MigrationSets {
		localSet := set

		r.Register(&CmdMigrateStatus{
			Set:        localSet,
			Connection: conn,
		})
		r.Register(&CmdMigrateDown{
			Set:        localSet,
			Connection: conn,
		})
		r.Register(&CmdMigrateUp{
			Set:        localSet,
			Connection: conn,
		})
	}

}

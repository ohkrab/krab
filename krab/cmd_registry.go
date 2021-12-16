package krab

// CmdRegistry is a list of registred commands.
type CmdRegistry []Cmd

// Register appends new command to registry.
func (r *CmdRegistry) Register(c Cmd) {
	*r = append(*r, c)
}

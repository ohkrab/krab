package krab

type Plugin interface {
	// Register is called once upon plugin initialization.
	// You can open connection here that will be closed later in Deregister.
	Register()

	// Deregister is called once at the end of plugin lifetime.
	Deregister()
}

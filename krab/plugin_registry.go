package krab

import "fmt"

// PluginRegistry data
type PluginRegistry struct {
	dir     string
	plugins []*Plugin
}

// NewPluginRegistry creates empty plugin registry for given path.
func NewPluginRegistry(dir string) *PluginRegistry {
	return &PluginRegistry{
		dir:     dir,
		plugins: make([]*Plugin, 0),
	}
}

// RegisterPlugins reads plugins file
func (pr *PluginRegistry) RegisterPlugins() {
	plugins := pr.readPlugins()
	for plugin := range plugins {
		fmt.Println("Loading ", plugin)
	}
}

func (pr *PluginRegistry) readPlugins() <-chan *Plugin {
	ch := make(chan *Plugin, 10)
	go func() {
		defer close(ch)
	}()
	return ch
}

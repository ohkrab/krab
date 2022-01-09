package krabwasm

import (
	"fmt"
	"sync"

	"github.com/wasmerio/wasmer-go/wasmer"
)

type Runtime struct {
	instances map[string]*webAssembly
	mutex     sync.RWMutex
	store     *wasmer.Store
}

func New() *Runtime {
	engine := wasmer.NewEngine()
	store := wasmer.NewStore(engine)

	return &Runtime{
		instances: map[string]*webAssembly{},
		store:     store,
	}
}

func (r *Runtime) LoadBytes(id string, b []byte) error {
	module, err := wasmer.NewModule(r.store, b)
	if err != nil {
		return fmt.Errorf("[wasm] Runtime failed to load bytes: %w", err)
	}
	importObj := wasmer.NewImportObject()
	// importObj.Register()
	instance, err := wasmer.NewInstance(module, importObj)

	if err != nil {
		return fmt.Errorf("[wasm] Runtime failed to instantiate module: %w", err)
	}

	r.mutex.Lock()
	r.instances[id] = &webAssembly{instance}
	r.mutex.Unlock()

	return nil
}

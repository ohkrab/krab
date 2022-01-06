package krabwasm

import (
	"github.com/wasmerio/wasmer-go/wasmer"
)

type webAssembly struct {
	instance *wasmer.Instance
}

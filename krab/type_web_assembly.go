package krab

import (
	"github.com/hashicorp/hcl2/hcl"
	"github.com/ohkrab/krab/krabhcl"
)

// WebAssembly represents web assembly config.
//
type WebAssembly struct {
	RefName string         `hcl:"ref_name,label"`
	Config  hcl.Attributes `hcl:"config"`
	File    string         `hcl:"file"`
}

func (w *WebAssembly) Addr() krabhcl.Addr {
	return krabhcl.Addr{Keyword: "wasm", Labels: []string{w.RefName}}
}

func (w *WebAssembly) Validate() error {
	return nil
}

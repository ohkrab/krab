package krab

import (
	"io"
)

// ToKCL converts DSL struct to Krab HCL.
type ToKCL interface {
	ToKCL(w io.StringWriter)
}

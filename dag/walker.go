package dag

import "github.com/ohkrab/krab/diagnostics"

type Walker struct {
	Callback WalkFunc
}

func (w *Walker) Walk() diagnostics.List {
	w.Callback(nil)
	return diagnostics.New()
}

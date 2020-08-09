package krab

import (
	"github.com/ohkrab/krab/dag"
	"github.com/ohkrab/krab/diagnostics"
)

type ContextGraphWalker struct {
	Context *Context
}

func (w *ContextGraphWalker) EvalContext() EvalContext {
	evaluator := &Evaluator{
		Config: w.Context.config,
	}
	ctx := DefaultEvalContext{
		Evaluator: evaluator,
	}
	return ctx
}

func (w *ContextGraphWalker) EnterVertex(dag.Vertex) {}

func (w *ContextGraphWalker) ExitVertex(dag.Vertex, diagnostics.List) {

}

func (w *ContextGraphWalker) EnterEvalTree(dag.Vertex, EvalNode) EvalNode {
	return nil
}

func (w *ContextGraphWalker) ExitEvalTree(dag.Vertex, interface{}, error) diagnostics.List {
	return nil
}

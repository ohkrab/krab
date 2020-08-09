package krab

import (
	"fmt"
	"log"

	"github.com/ohkrab/krab/configs"
	"github.com/ohkrab/krab/dag"
	"github.com/ohkrab/krab/diagnostics"
)

type GraphTransformer interface {
	Transform(*Graph) error
}

type Graph struct {
	graph  *dag.Graph
	module *configs.Module
}

type GraphBuilder struct {
	Steps    []GraphTransformer
	Validate bool
}

type GraphWalker interface {
	EvalContext() EvalContext
	EnterVertex(dag.Vertex)
	ExitVertex(dag.Vertex, diagnostics.List)
	EnterEvalTree(dag.Vertex, EvalNode) EvalNode
	ExitEvalTree(dag.Vertex, interface{}, error) diagnostics.List
}

func (b *GraphBuilder) Build(mod *configs.Module) (*Graph, diagnostics.List) {
	g := &Graph{
		module: mod,
		graph:  dag.New(),
	}
	diags := diagnostics.New()

	for _, step := range b.Steps {
		log.Printf("[TRACE] Executing graph transform %T", step)
		if err := step.Transform(g); err != nil {
			diags.Append(err)
		}
	}

	return g, diags
}

func (g *Graph) Walk(walker GraphWalker) diagnostics.List {
	ctx := walker.EvalContext()

	walkFn := func(v dag.Vertex) (diags diagnostics.List) {
		log.Printf("[TRACE] vertex %q: starting visit (%T)", v.VertexID(), v)

		defer func() {
			log.Printf("[TRACE] vertex %q: visit complete", v.VertexID())
		}()

		walker.EnterVertex(v)
		defer walker.ExitVertex(v, diags)

		// If the node is eval-able, then evaluate it.
		if ev, ok := v.(GraphNodeEvalable); ok {
			tree := ev.EvalTree()
			if tree == nil {
				panic(fmt.Sprintf("%q (%T): nil eval tree", v.VertexID(), v))
			}

			// Allow the walker to change our tree if needed. Eval,
			// then callback with the output.
			log.Printf("[TRACE] vertex %q: evaluating", v.VertexID())

			tree = walker.EnterEvalTree(v, tree)
			output, err := tree.Eval(ctx)
			diags = diags.Append(walker.ExitEvalTree(v, output, err))
			if diags.HasErrors() {
				return
			}
		}

		return
	}

	return g.graph.Walk(walkFn)
}

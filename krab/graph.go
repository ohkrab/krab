package krab

import (
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
type EvalNode interface{}

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
	var walkFn dag.WalkFunc
	// ctx := walker.EvalContext()
	return g.graph.Walk(walkFn)
}

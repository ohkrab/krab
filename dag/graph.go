package dag

import (
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/ohkrab/krab/diagnostics"
)

type Vertex interface {
	VertexID() string
}

type WalkFunc func(Vertex) diagnostics.List

type Graph struct {
	data    map[string]Vertex
	adjList map[string]*treeset.Set
}

func New() *Graph {
	return &Graph{
		data:    make(map[string]Vertex, 10),
		adjList: make(map[string]*treeset.Set, 10),
	}
}

func (g *Graph) HasVertex(v string) bool {
	_, ok := g.adjList[v]
	return ok
}

func (g *Graph) AddVertex(v Vertex) {
	if !g.HasVertex(v.VertexID()) {
		g.adjList[v.VertexID()] = treeset.NewWithStringComparator()
		g.data[v.VertexID()] = v
	}
}

func (g *Graph) AddEdge(from, to string) {
	if g.HasVertex(from) && g.HasVertex(to) {
		if !g.adjList[from].Contains(to) {
			g.adjList[from].Add(to)
		}
	}
}

func (g *Graph) Walk(walkFn WalkFunc) diagnostics.List {
	w := Walker{Callback: walkFn, Graph: g}
	return w.Walk()
}

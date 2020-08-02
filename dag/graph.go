package dag

import (
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/ohkrab/krab/diagnostics"
)

type Vertex interface{}

type WalkFunc func(Vertex) diagnostics.List

type Graph struct {
	adjList map[Vertex]*hashset.Set
}

func New() *Graph {
	return &Graph{
		adjList: make(map[Vertex]*hashset.Set, 0),
	}
}

func (g *Graph) HasVertex(v Vertex) bool {
	_, ok := g.adjList[v]
	return ok
}

func (g *Graph) AddVertex(v Vertex) {
	if !g.HasVertex(v) {
		g.adjList[v] = hashset.New()
	}
}

func (g *Graph) AddEdge(from, to Vertex) {
	if g.HasVertex(from) && g.HasVertex(to) {
		if !g.adjList[from].Contains(to) {
			g.adjList[from].Add(to)
		}
	}
}

func (g *Graph) Walk(walkFn WalkFunc) diagnostics.List {
	w := &Walker{Callback: walkFn}
	return w.Walk()
}

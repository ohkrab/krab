package dag

import (
	"fmt"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/ohkrab/krab/diagnostics"
)

type Vertex interface {
	VertexID() string
}

type WalkFunc func(Vertex) diagnostics.List

type Graph struct {
	data    map[string]Vertex
	adjList *treemap.Map
}

func New() *Graph {
	return &Graph{
		data:    make(map[string]Vertex, 10),
		adjList: treemap.NewWithStringComparator(),
	}
}

func (g *Graph) HasVertex(v string) bool {
	_, ok := g.adjList.Get(v)
	return ok
}

func (g *Graph) AddVertex(v Vertex) {
	if !g.HasVertex(v.VertexID()) {
		g.adjList.Put(
			v.VertexID(),
			treeset.NewWithStringComparator(),
		)
		g.data[v.VertexID()] = v
	}
}

func (g *Graph) AddEdge(from, to string) {
	if g.HasVertex(from) && g.HasVertex(to) {
		set := g.adjecentSet(from)
		if !set.Contains(to) {
			set.Add(to)
		}
	}
}

func (g *Graph) Walk(walkFn WalkFunc) diagnostics.List {
	w := Walker{Callback: walkFn, Graph: g}
	return w.Walk()
}

func (g *Graph) eachVertex(eachFn func(v string)) {
	g.adjList.Each(func(index interface{}, value interface{}) {
		if v, ok := index.(string); ok {
			eachFn(v)
		} else {
			panic("Graph: Vertex key is not a `string`")
		}
	})
}

func (g *Graph) adjecentSet(v string) *treeset.Set {
	if list, found := g.adjList.Get(v); found {
		if set, ok := list.(*treeset.Set); ok {
			return set
		}
	}

	panic(fmt.Sprintf("Graph: adjecentSet should always return `*treeset.Set`"))
}

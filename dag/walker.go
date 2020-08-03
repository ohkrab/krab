package dag

import (
	"github.com/ohkrab/krab/diagnostics"
)

type Walker struct {
	*Graph
	Callback WalkFunc
}

func (w *Walker) Walk() diagnostics.List {
	diags := diagnostics.New()
	visited := make(map[string]bool, 0)

	w.Graph.eachVertex(func(v string) {
		if !visited[v] {
			w.recursiveWalk(v, visited, diags)
		}
	})

	return diags
}

func (w *Walker) recursiveWalk(startingVertex string, visited map[string]bool, diags diagnostics.List) {
	w.Graph.eachAdjecentVertex(startingVertex, func(v string) {
		if !visited[v] {
			w.recursiveWalk(v, visited, diags)
		}
	})

	visited[startingVertex] = true
	w.Callback(w.Graph.data[startingVertex])
}

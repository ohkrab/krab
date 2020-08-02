package dag

import (
	"errors"

	"github.com/ohkrab/krab/diagnostics"
)

type Walker struct {
	*Graph
	Callback WalkFunc
}

func (w *Walker) Walk() diagnostics.List {
	diags := diagnostics.New()
	visited := make(map[string]bool, 0)

	for v := range w.Graph.adjList {
		if !visited[v] {
			w.recursiveWalk(v, visited, diags)
		}
	}

	return diags
}

func (w *Walker) recursiveWalk(startingVertex string, visited map[string]bool, diags diagnostics.List) {
	adjList := w.Graph.adjList[startingVertex]
	it := adjList.Iterator()

	for it.Next() {
		v, ok := it.Value().(string)
		if ok {
			if !visited[v] {
				w.recursiveWalk(v, visited, diags)
			}
		} else {
			diags.Append(errors.New("Cannot fetch vertex ID in recursive walk"))
		}
	}

	visited[startingVertex] = true
	w.Callback(w.Graph.data[startingVertex])
}

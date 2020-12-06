package dag

import (
	"fmt"
	"strings"
	"testing"

	"github.com/franela/goblin"
	"github.com/ohkrab/krab/diagnostics"
)

type MockVertex struct {
	ID int
}

func (mv MockVertex) VertexID() string {
	return fmt.Sprint(mv.ID)
}

func Test_Walker(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Walker#Walk", func() {
		g.It("Traverse graph using DFS", func() {
			graph := New()
			graph.AddVertex(MockVertex{ID: 1})
			graph.AddVertex(MockVertex{ID: 2})
			graph.AddVertex(MockVertex{ID: 3})
			graph.AddVertex(MockVertex{ID: 4})
			graph.AddVertex(MockVertex{ID: 5})
			graph.AddVertex(MockVertex{ID: 6})
			//    (1)
			//   /   \
			//  (2)  (3)
			//  |  \ /
			// (4) (5)
			//  |
			// (6)
			graph.AddEdge("1", "2")
			graph.AddEdge("1", "3")
			graph.AddEdge("2", "4")
			graph.AddEdge("2", "5")
			graph.AddEdge("3", "5")
			graph.AddEdge("4", "6")

			visited := make([]string, 0, 6)
			graph.Walk(func(v Vertex) diagnostics.List {
				visited = append(visited, v.VertexID())
				return nil
			})

			walkedPath := strings.Join(visited, ",")

			g.Assert(walkedPath).Eql("6,4,5,2,3,1")
		})
	})
}

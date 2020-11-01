package configs

import (
	"fmt"
	"testing"

	. "github.com/franela/goblin"
)

func Test_ParserConnections(t *testing.T) {
	g := Goblin(t)
	g.Describe("Parser#Connections", func() {
		g.It("Should parse connections", func() {
			parser := NewParser()
			c, diags := parser.LoadConfigDir("./../test/fixtures")

			g.Assert(len(diags)).Eql(0)

			fmt.Println(c.Module.Connections)
		})
	})
}

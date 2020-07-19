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
			file, diags := parser.LoadConfigFile("./../test/fixtures/connections.krab")

			for _, c := range file.Connections {
				fmt.Println(c.Name, ":", c.UriVal)
			}

			fmt.Println(diags)
		})
	})
}

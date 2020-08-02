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
			config, diags := parser.LoadConfigDir("./../test/fixtures")

			for addr, conn := range config.Module.Connections {
				fmt.Println(addr, ":", conn)
			}

			fmt.Println(diags)
		})
	})
}

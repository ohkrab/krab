package krab

import (
	"fmt"
	"testing"

	"github.com/franela/goblin"
)

func Test_ContextEval(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Context#Eval", func() {
		g.It("Should evaluate graph", func() {
			ctx, diags := Load("./../test/fixtures")

			fmt.Println(diags)

			ctx.Eval()
		})
	})
}

package krab

import (
	"testing"

	"github.com/franela/goblin"
)

func TestValidators(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("ValidateRefName", func() {
		g.It("Allow alphanumeric and underscore", func() {
			g.Assert(ValidateRefName("valid_ref")).IsNil()
			g.Assert(ValidateRefName("valid_123")).IsNil()
			g.Assert(ValidateRefName("ValidRef")).IsNil()
			g.Assert(ValidateRefName("___")).IsNil()
		})

		g.It("Cannot start with number", func() {
			g.Assert(ValidateRefName("123")).IsNotNil()
			g.Assert(ValidateRefName("123_abc")).IsNotNil()
		})

		g.It("Cannot be empty", func() {
			g.Assert(ValidateRefName("")).IsNotNil()
		})

		g.It("No other separators", func() {
			g.Assert(ValidateRefName("abc-def")).IsNotNil()
			g.Assert(ValidateRefName("abc def")).IsNotNil()
		})
	})

	g.Describe("ValidateStringNonEmpty", func() {
		g.It("Lenght must be > 0", func() {
			g.Assert(ValidateStringNonEmpty("field", "a")).IsNil()
			g.Assert(ValidateStringNonEmpty("field", "")).IsNotNil()
		})
	})
}

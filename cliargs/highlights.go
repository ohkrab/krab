package cliargs

import (
	"fmt"

	"github.com/wzshiming/ctc"
)

func Highlight(color ctc.Color, s string, colorize bool) string {
	if colorize {
		return fmt.Sprint(color, s, ctc.Reset)
	}
	return fmt.Sprint(ctc.Reset, s, ctc.Reset)
}

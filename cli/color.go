package cli

import (
	"fmt"

	"github.com/wzshiming/ctc"
)

// Red colorizes output.
func Red(s string) string {
	return fmt.Sprint(ctc.ForegroundRed, s, ctc.Reset)
}

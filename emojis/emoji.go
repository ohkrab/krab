package emojis

import (
	"fmt"

	"github.com/wzshiming/ctc"
)

// CheckMarkColor ✔ emoji
func CheckMarkColor(color ctc.Color) string {
	return fmt.Sprintf("%s%c%s ", color, '\u2714', ctc.Reset)
}

// CheckMark green check mark
func CheckMark() string {
	return CheckMarkColor(ctc.ForegroundGreen)
}

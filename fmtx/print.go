package fmtx

import (
	"fmt"
	"os"

	"github.com/wzshiming/ctc"
)

func WriteError(format string, a ...any) {
	fmt.Fprintf(os.Stderr, fmt.Sprint(ctc.ForegroundRed, format, ctc.Reset, "\n"), a...)
}

func WriteSuccess(format string, a ...any) {
	fmt.Fprintf(os.Stdout, fmt.Sprint(ctc.ForegroundGreen, format, ctc.Reset, "\n"), a...)
}

func WriteInfo(format string, a ...any) {
	fmt.Fprintf(os.Stdout, fmt.Sprint(ctc.ForegroundCyan, format, ctc.Reset, "\n"), a...)
}

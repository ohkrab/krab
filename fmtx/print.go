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

func WriteLine(format string, a ...any) {
	fmt.Fprintf(os.Stdout, format + "\n", a...)
}

func ColoredBlockDanger(format string, a ...any) string {
	return fmt.Sprintf(
		fmt.Sprint(
			ctc.BackgroundBrightRed | ctc.ForegroundBlack,
			format,
			ctc.Reset,
		),
		a...,
	)
}

func ColoredBlockSuccess(format string, a ...any) string {
	return fmt.Sprintf(
		fmt.Sprint(
			ctc.BackgroundBrightGreen | ctc.ForegroundBlack,
			format,
			ctc.Reset,
		),
		a...,
	)
}

func ColoredBlockWarning(format string, a ...any) string {
	return fmt.Sprintf(
		fmt.Sprint(
			ctc.BackgroundBrightYellow | ctc.ForegroundBlack,
			format,
			ctc.Reset,
		),
		a...,
	)
}

func Danger(format string, a ...any) string {
    return fmt.Sprintf(
        fmt.Sprint(ctc.ForegroundRed, format, ctc.Reset),
        a...,
    )
}

func Success(format string, a ...any) string {
    return fmt.Sprintf(
        fmt.Sprint(ctc.ForegroundGreen, format, ctc.Reset),
        a...,
    )
}

func Warning(format string, a ...any) string {
    return fmt.Sprintf(
        fmt.Sprint(ctc.ForegroundYellow, format, ctc.Reset),
        a...,
    )
}


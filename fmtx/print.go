package fmtx

import (
	"fmt"
	"io"

	"github.com/wzshiming/ctc"
)

func (l Logger) WriteError(format string, a ...any) {
	FWriteError(l.Stderr, format, a...)
}

func (l Logger) WriteSuccess(format string, a ...any) {
	FWriteSuccess(l.Stdout, format, a...)
}

func (l Logger) WriteInfo(format string, a ...any) {
	FWriteInfo(l.Stdout, format, a...)
}

func (l Logger) WriteLine(format string, a ...any) {
	FWriteLine(l.Stdout, format, a...)
}

func FWriteError(w io.Writer, format string, a ...any) {
	fmt.Fprintf(w, fmt.Sprint(ctc.ForegroundRed, format, ctc.Reset, "\n"), a...)
}

func FWriteSuccess(w io.Writer, format string, a ...any) {
	fmt.Fprintf(w, fmt.Sprint(ctc.ForegroundGreen, format, ctc.Reset, "\n"), a...)
}

func FWriteInfo(w io.Writer, format string, a ...any) {
	fmt.Fprintf(w, fmt.Sprint(ctc.ForegroundCyan, format, ctc.Reset, "\n"), a...)
}

func FWriteLine(w io.Writer, format string, a ...any) {
	fmt.Fprintf(w, format+"\n", a...)
}

func ColoredBlockDanger(format string, a ...any) string {
	return fmt.Sprintf(
		fmt.Sprint(
			ctc.BackgroundBrightRed|ctc.ForegroundBlack,
			format,
			ctc.Reset,
		),
		a...,
	)
}

func ColoredBlockSuccess(format string, a ...any) string {
	return fmt.Sprintf(
		fmt.Sprint(
			ctc.BackgroundBrightGreen|ctc.ForegroundBlack,
			format,
			ctc.Reset,
		),
		a...,
	)
}

func ColoredBlockWarning(format string, a ...any) string {
	return fmt.Sprintf(
		fmt.Sprint(
			ctc.BackgroundBrightYellow|ctc.ForegroundBlack,
			format,
			ctc.Reset,
		),
		a...,
	)
}

func Danger(format string, a ...any) string {
	return fmt.Sprintf(
		fmt.Sprint(ctc.ForegroundBrightRed, format, ctc.Reset),
		a...,
	)
}

func Success(format string, a ...any) string {
	return fmt.Sprintf(
		fmt.Sprint(ctc.ForegroundBrightGreen, format, ctc.Reset),
		a...,
	)
}

func Warning(format string, a ...any) string {
	return fmt.Sprintf(
		fmt.Sprint(ctc.ForegroundBrightYellow, format, ctc.Reset),
		a...,
	)
}

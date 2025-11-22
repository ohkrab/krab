package fmtx

import (
	"io"
	"os"
)

type Logger struct {
	Stdout io.Writer
	Stderr io.Writer
}

func New(out io.Writer, err io.Writer) *Logger {
	l := &Logger{
		Stdout: out,
		Stderr: err,
	}

	return l
}

func Default() *Logger {
	return New(os.Stdout, os.Stderr)
}

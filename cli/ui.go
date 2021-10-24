package cli

import (
	"io"
	"os"

	mcli "github.com/mitchellh/cli"
)

// UI implements the mitchellh/cli.Ui interface.
type UI interface {
	mcli.Ui
}

func New(errorWriter io.Writer, writer io.Writer) UI {
	ui := &mcli.ColoredUi{
		Ui:         &mcli.BasicUi{ErrorWriter: errorWriter, Writer: writer},
		WarnColor:  mcli.UiColorYellow,
		ErrorColor: mcli.UiColorRed,
		InfoColor:  mcli.UiColorGreen,
	}

	return ui
}

func DefaultUI() UI {
	return New(os.Stderr, os.Stdout)
}

func NullUI() UI {
	ui := &mcli.BasicUi{ErrorWriter: io.Discard, Writer: io.Discard}

	return ui
}

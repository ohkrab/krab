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

func DefaultUI() UI {
	ui := &mcli.ColoredUi{
		Ui:         &mcli.BasicUi{ErrorWriter: os.Stderr, Writer: os.Stdout},
		WarnColor:  mcli.UiColorYellow,
		ErrorColor: mcli.UiColorRed,
		InfoColor:  mcli.UiColorGreen,
	}

	return ui
}

func NullUI() UI {
	ui := &mcli.BasicUi{ErrorWriter: io.Discard, Writer: io.Discard}

	return ui
}

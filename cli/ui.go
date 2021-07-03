package cli

import (
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

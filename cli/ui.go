package cli

import (
	"io"
	"os"

	mcli "github.com/mitchellh/cli"
)

// UI implements the mitchellh/cli.Ui interface.
type UI interface {
	mcli.Ui
	Stdout() io.Writer
	Stderr() io.Writer
}

type BasicUI struct {
	mcli.BasicUi
}

func (u *BasicUI) Stdout() io.Writer {
	return u.BasicUi.Writer
}

func (u *BasicUI) Stderr() io.Writer {
	return u.BasicUi.ErrorWriter
}

func DefaultUI() *mcli.ColoredUi {
	basicUI := &BasicUI{
		BasicUi: mcli.BasicUi{ErrorWriter: os.Stderr, Writer: os.Stdout},
	}
	ui := &mcli.ColoredUi{
		Ui:         basicUI,
		WarnColor:  mcli.UiColorYellow,
		ErrorColor: mcli.UiColorRed,
	}

	return ui
}
